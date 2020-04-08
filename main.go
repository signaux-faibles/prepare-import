package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"io/ioutil"
)

// Implementation of the prepare-import command.
func main() {
	var path = flag.String("path", ".", "Chemin d'accès aux fichiers données")
	flag.Parse()
	adminObject, err := PrepareImport(*path)
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

// FilesProperty represents the "files" property of an Admin object.
type FilesProperty map[string][]string

// DataFile represents a Data File to be imported, and allows to determine its type and name.
type DataFile interface {
	GetFilename() string    // the name as it will be stored in Admin
	DetectFileType() string // returns the type of that file (e.g. "debit")
}

// SimpleDataFile is a DataFile which type can be determined without requiring a metadata file (e.g. well-named csv file).
type SimpleDataFile struct {
	filename string
}

func (dataFile SimpleDataFile) DetectFileType() string {
	return ExtractFileTypeFromFilename(dataFile.filename)
}

func (dataFile SimpleDataFile) GetFilename() string {
	return dataFile.filename
}

// UploadedDataFile is a DataFile which type can be determined thanks to a metadata file (e.g. bin+info files).
type UploadedDataFile struct {
	filename string
	path     string
}

func (dataFile UploadedDataFile) DetectFileType() string {
	metaFilepath := filepath.Join(dataFile.path, strings.Replace(dataFile.filename, ".bin", ".info", 1))
	fileinfo := LoadMetadata(metaFilepath)
	return ExtractFileTypeFromMetadata(metaFilepath, fileinfo) // e.g. "Sigfaible_debits.csv"
}

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
func PrepareImport(pathname string) (AdminObject, error) {
	filenames, err := ReadFilenames(pathname)
	if err != nil {
		return nil, err
	}
	augmentedFiles := []DataFile{}
	for _, file := range filenames {
		augmentedFiles = append(augmentedFiles, AugmentDataFile(file, pathname))
	}
	return PurePrepareImport(augmentedFiles)
}

// PurePrepareImport populates an AdminObject, given a list of data files.
func PurePrepareImport(augmentedFilenames []DataFile) (AdminObject, error) {
	filesProperty, unsupportedFiles := PopulateFilesProperty(augmentedFilenames)
	var err error
	if len(unsupportedFiles) > 0 {
		err = UnsupportedFilesError{unsupportedFiles}
	}
	var completeTypes = []string{}
	for _, typeName := range defaultCompleteTypes {
		if _, ok := filesProperty[typeName]; ok {
			completeTypes = append(completeTypes, typeName)
		}
	}
	return AdminObject{"files": filesProperty, "complete_types": completeTypes}, err
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
