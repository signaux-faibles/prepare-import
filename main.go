package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"io/ioutil"
)

func main() {
	var path = flag.String("path", ".", "Chemin d'accès aux fichiers données")
	flag.Parse()
	adminObject, err := PrepareImport(*path)
	if err != nil {
		log.Fatal(err)
	}
	json, err := json.MarshalIndent(adminObject, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(json))
}

type AdminObject map[string]interface{}

type DataFile interface {
	GetFilename() string // the name as it will be stored in Admin
	DetectFileType() string
}

type SimpleDataFile struct {
	filename string
}

func (dataFile SimpleDataFile) DetectFileType() string {
	return ExtractFileTypeFromFilename(dataFile.filename)
}

func (dataFile SimpleDataFile) GetFilename() string {
	return dataFile.filename
}

type UploadedDataFile struct {
	filename string
	path     string
}

func (dataFile UploadedDataFile) DetectFileType() string {
	metaFilepath := filepath.Join(dataFile.path, strings.Replace(dataFile.filename, ".bin", ".info", 1))
	fileinfo, err := LoadMetadata(metaFilepath)
	if err != nil {
		log.Fatal(err)
	}
	filetype := ExtractFileTypeFromMetadata(metaFilepath, fileinfo)
	return filetype // e.g. "Sigfaible_debits.csv"
}

func (dataFile UploadedDataFile) GetFilename() string {
	return dataFile.filename
}

func PrepareImport(pathname string) (AdminObject, error) {
	filenames, err := ReadFilenames(pathname)
	if err != nil {
		return nil, err
	}
	augmentedFiles := []DataFile{}
	for _, file := range filenames {
		var filename DataFile
		if strings.HasSuffix(file, ".bin") {
			filename = UploadedDataFile{file, pathname}
		} else {
			filename = SimpleDataFile{file}
		}
		augmentedFiles = append(augmentedFiles, filename)
	}
	return PurePrepareImport(augmentedFiles), nil
}

func PurePrepareImport(augmentedFilenames []DataFile) AdminObject {
	filesProperty := PopulateFilesProperty(augmentedFilenames)
	return AdminObject{"files": filesProperty}
}

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

type FilesProperty map[string][]string

func LoadMetadata(filepath string) (UploadedFileMeta, error) {

	// read file
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return UploadedFileMeta{}, err
	}

	// unmarshall data from json
	var uploadedFileMeta UploadedFileMeta
	err = json.Unmarshal(data, &uploadedFileMeta)
	if err != nil {
		return UploadedFileMeta{}, err
	}

	return uploadedFileMeta, nil
}

func PopulateFilesProperty(filenames []DataFile) FilesProperty {
	filesProperty := FilesProperty{}
	for _, filename := range filenames {
		filetype := filename.DetectFileType()

		if filetype == "" {
			// Unsupported file
			continue
		}
		if _, exists := filesProperty[filetype]; !exists {
			filesProperty[filetype] = []string{}
		}
		filesProperty[filetype] = append(filesProperty[filetype], filename.GetFilename())
	}
	return filesProperty
}
