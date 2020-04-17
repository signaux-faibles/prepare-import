package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"io/ioutil"
)

// Implementation of the prepare-import command.
func main() {
	var path = flag.String("path", ".", "Chemin d'accès aux fichiers données")
	var batchKey = flag.String(
		"batch",
		"",
		"Clé du batch à importer au format AAMM (année + mois + suffixe optionnel)\n"+
			"Exemple: 1802_1",
	)
	var dateFinEffectif = flag.String(
		"date-fin-effectif",
		"",
		"Date de fin des données \"effectif\" fournies, au format AAAA-MM-JJ (année + mois + jour)\n"+
			"Exemple: 2014-01-01",
	)
	flag.Parse()
	validBatchKey, err := BatchKey(*batchKey)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error()+"\n\nUsage:")
		flag.PrintDefaults()
		os.Exit(1)
	}
	validDateFinEffectif, err := DateFinEffectif(*dateFinEffectif)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error()+"\n\nUsage:")
		flag.PrintDefaults()
		os.Exit(1)
	}
	adminObject, err := PrepareImport(*path, validBatchKey, validDateFinEffectif)
	if _, ok := err.(UnsupportedFilesError); ok {
		fmt.Fprintln(os.Stderr, err.Error())
	} else if err != nil {
		log.Fatal(err) // will print in the error output stream and exit
	}
	json, err := json.MarshalIndent(adminObject, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(json))
}

type batchKeyType string

// BatchKey instanciates a batchKeyType after validating its value.
func BatchKey(key string) (batchKeyType, error) {
	var isValidBatchKey = regexp.MustCompile(`^[0-9]{4}`)
	if !isValidBatchKey.MatchString(key) {
		return batchKeyType(""), errors.New("la clé du batch doit respecter le format requis AAMM")
	}
	return batchKeyType(key), nil
}

type dateFinEffectifType string

// DateFinEffectif instanciates a dateFinEffectifType after validating its value.
func DateFinEffectif(date string) (dateFinEffectifType, error) {
	var isDateFinEffectif = regexp.MustCompile(`^[0-9]{4}-[0-1][0-9]-[0-3][0-9]$`)
	if !isDateFinEffectif.MatchString(date) {
		return dateFinEffectifType(""), errors.New("la date-fin-effectif doit respecter le format requis AAAA-MM-JJ")
	}
	return dateFinEffectifType(date), nil
}

// AdminObject represents a document going to be stored in the Admin db collection.
type AdminObject map[string]interface{}

// IDProperty represents the "_id" property of an Admin object.
type IDProperty struct {
	Key  batchKeyType `json:"key"`
	Type string       `json:"type"`
}

// FilesProperty represents the "files" property of an Admin object.
type FilesProperty map[ValidFileType][]string

// MongoDate represents a date that can be serialized for MongoDB.
type MongoDate struct {
	date string
}

// MarshalJSON will be called when serializing a date for MongoDB.
func (mongoDate MongoDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{"$date": mongoDate.date})
}

// UnmarshalJSON will parse a MongoDate from a MongoDB date object.
func (mongoDate *MongoDate) UnmarshalJSON(data []byte) error {
	var dateObj map[string]string
	if err := json.Unmarshal(data, &dateObj); err != nil {
		return err
	}
	mongoDate.date = dateObj["$date"]
	return nil
}

// ParamProperty represents the "param" property of an Admin object.
type ParamProperty struct {
	DateDebut       MongoDate `json:"date_debut"`
	DateFin         MongoDate `json:"date_fin"`
	DateFinEffectif MongoDate `json:"date_fin_effectif"`
}

// DataFile represents a Data File to be imported, and allows to determine its type and name.
type DataFile interface {
	GetFilename() string           // the name as it will be stored in Admin
	DetectFileType() ValidFileType // returns the type of that file (e.g. DEBIT)
}

// SimpleDataFile is a DataFile which type can be determined without requiring a metadata file (e.g. well-named csv file).
type SimpleDataFile struct {
	filename string
}

// DetectFileType returns the type of that file (e.g. DEBIT).
func (dataFile SimpleDataFile) DetectFileType() ValidFileType {
	return ExtractFileTypeFromFilename(dataFile.filename)
}

// GetFilename returns the name as it will be stored in Admin.
func (dataFile SimpleDataFile) GetFilename() string {
	return dataFile.filename
}

// UploadedDataFile is a DataFile which type can be determined thanks to a metadata file (e.g. bin+info files).
type UploadedDataFile struct {
	filename string
	path     string
}

// DetectFileType returns the type of that file (e.g. DEBIT).
func (dataFile UploadedDataFile) DetectFileType() ValidFileType {
	metaFilepath := filepath.Join(dataFile.path, strings.Replace(dataFile.filename, ".bin", ".info", 1))
	fileinfo := LoadMetadata(metaFilepath)
	return ExtractFileTypeFromMetadata(metaFilepath, fileinfo) // e.g. "Sigfaible_debits.csv"
}

// GetFilename returns the name as it will be stored in Admin.
func (dataFile UploadedDataFile) GetFilename() string {
	return dataFile.filename
}

// UnsupportedFilesError is an Error object that lists files that were not supported.
type UnsupportedFilesError struct {
	UnsupportedFiles []string
}

func (err UnsupportedFilesError) Error() string {
	return "unsupported: " + strings.Join(err.UnsupportedFiles, ", ")
}

// ReadFilenames returns the name of files found at the provided path.
func ReadFilenames(path string) ([]string, error) {
	var files []string
	fileInfo, err := ioutil.ReadDir(path)
	if err != nil {
		return files, err
	}
	for _, file := range fileInfo {
		files = append(files, file.Name())
	}
	return files, nil
}

// AugmentDataFile returns a SimpleDataFile or UploadedDataFile (if metadata had to be loaded).
func AugmentDataFile(file string, pathname string) DataFile {
	if strings.HasSuffix(file, ".bin") {
		return UploadedDataFile{file, pathname}
	}
	return SimpleDataFile{file}
}

// PrepareImport generates an Admin object from files found at given pathname of the file system.
func PrepareImport(pathname string, batchKey batchKeyType, dateFinEffectif dateFinEffectifType) (AdminObject, error) {
	filenames, err := ReadFilenames(pathname)
	if err != nil {
		return nil, err
	}
	augmentedFiles := []DataFile{}
	for _, file := range filenames {
		augmentedFiles = append(augmentedFiles, AugmentDataFile(file, pathname))
	}
	return PopulateAdminObject(augmentedFiles, batchKey, dateFinEffectif)
}

// PopulateAdminObject populates an AdminObject, given a list of data files.
func PopulateAdminObject(augmentedFilenames []DataFile, batchKey batchKeyType, dateFinEffectif dateFinEffectifType) (AdminObject, error) {

	filesProperty, unsupportedFiles := PopulateFilesProperty(augmentedFilenames)
	var err error
	if len(unsupportedFiles) > 0 {
		err = UnsupportedFilesError{unsupportedFiles}
	}
	var completeTypes = []ValidFileType{}
	for _, typeName := range defaultCompleteTypes {
		if _, ok := filesProperty[typeName]; ok {
			completeTypes = append(completeTypes, typeName)
		}
	}
	// { "date_debut" : { "$date" : "2014-01-01T00:00:00.000+0000" }, "date_fin" : { "$date" : "2018-12-01T00:00:00.000+0000" }, "date_fin_effectif" : { "$date" : "2018-06-01T00:00:00.000+0000" } }

	paramProperty := ParamProperty{
		DateDebut:       MongoDate{"2014-01-01T00:00:00.000+0000"},
		DateFin:         MongoDate{"20" + string(batchKey)[0:2] + "-" + string(batchKey)[2:4] + "-01T00:00:00.000+0000"},
		DateFinEffectif: MongoDate{string(dateFinEffectif) + "T00:00:00.000+0000"},
	}

	return AdminObject{
		"_id":            IDProperty{batchKey, "batch"},
		"files":          filesProperty,
		"complete_types": completeTypes,
		"param":          paramProperty,
	}, err
}

// LoadMetadata returns the metadata of a .bin file, by reading the given .info file.
func LoadMetadata(filepath string) UploadedFileMeta {

	// read file
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	// unmarshall data from json
	var uploadedFileMeta UploadedFileMeta
	err = json.Unmarshal(data, &uploadedFileMeta)
	if err != nil {
		panic(err)
	}

	return uploadedFileMeta
}

// PopulateFilesProperty populates the "files" property of an Admin object, given a list of Data files.
func PopulateFilesProperty(filenames []DataFile) (FilesProperty, []string) {
	filesProperty := FilesProperty{}
	unsupportedFiles := []string{}
	for _, filename := range filenames {
		filetype := filename.DetectFileType()

		if filetype == "" {
			if !strings.HasSuffix(filename.GetFilename(), ".info") {
				unsupportedFiles = append(unsupportedFiles, filename.GetFilename())
			}
			continue
		}
		if _, exists := filesProperty[filetype]; !exists {
			filesProperty[filetype] = []string{}
		}
		filesProperty[filetype] = append(filesProperty[filetype], filename.GetFilename())
	}
	return filesProperty, unsupportedFiles
}
