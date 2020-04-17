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
	flag.Parse()
	validBatchKey, err := NewBatchKey(*batchKey)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error()+"\n\nUsage:")
		flag.PrintDefaults()
		os.Exit(1)
	}
	adminObject, err := PrepareImport(*path, validBatchKey)
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

// AdminObject represents a document going to be stored in the Admin db collection.
type AdminObject map[string]interface{}

// IDProperty represents the "_id" property of an Admin object.
type IDProperty struct {
	Key  BatchKey `json:"key"`
	Type string   `json:"type"`
}

// FilesProperty represents the "files" property of an Admin object.
type FilesProperty map[ValidFileType][]string

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
func PrepareImport(pathname string, batchKey BatchKey) (AdminObject, error) {
	filenames, err := ReadFilenames(pathname)
	if err != nil {
		return nil, err
	}
	augmentedFiles := []DataFile{}
	for _, file := range filenames {
		augmentedFiles = append(augmentedFiles, AugmentDataFile(file, pathname))
	}
	return PopulateAdminObject(augmentedFiles, batchKey)
}

type batchKeyType string

func (b batchKeyType) String() string {
	return string(b)
}

// BatchKey represents a valid batch key.
type BatchKey interface {
	String() string
}

// NewBatchKey constructs a valid batch key.
func NewBatchKey(key string) (BatchKey, error) {
	var isValidBatchKey = regexp.MustCompile(`^[0-9]{4}`)
	if !isValidBatchKey.MatchString(key) {
		return batchKeyType(""), errors.New("la clé du batch doit respecter le format requis AAMM")
	}
	return batchKeyType(key), nil
}

// PopulateAdminObject populates an AdminObject, given a list of data files.
func PopulateAdminObject(augmentedFilenames []DataFile, batchKey BatchKey) (AdminObject, error) {

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

	paramProperty := map[string]map[string]string{
		"date_debut": map[string]string{"$date": "2014-01-01T00:00:00.000+0000"},
		"date_fin":   map[string]string{"$date": "20" + batchKey.String()[0:2] + "-" + batchKey.String()[2:4] + "-01T00:00:00.000+0000"},
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
