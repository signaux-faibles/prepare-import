package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	// flags
	var path = flag.String("path", ".", "Chemin d'accès aux fichiers données")
	flag.Parse()
	adminObject, err := PrepareImport(*path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(adminObject)
	// serialization
}

type AdminObject map[string]interface{}

func PrepareImport(pathname string) (AdminObject, error) {
	filenames, err := ReadFilenames(pathname)
	if err != nil {
		return nil, err
	}
	return PurePrepareImport(filenames)
}

func PurePrepareImport(filenames []string) (AdminObject, error) {
	fileProperty := PopulateFilesProperty(filenames)
	return AdminObject{"files": fileProperty}, nil
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

type FileProperty map[string][]string

func PopulateFilesProperty(filenames []string) FileProperty {
	fileProperty := FileProperty{
		// "effectif": []string{"coucou"},
		// "debit":    []string{},
	}
	for _, filename := range filenames {
		filetype, _ := GetFileType(filename)
    if filetype == "" {
      // Unsupported file
      continue
    }
		if _, exists := fileProperty[filetype]; !exists {
			fileProperty[filetype] = []string{}
		}
		fileProperty[filetype] = append(fileProperty[filetype], filename)
	}
	return fileProperty
}

func GetFileType(filename string) (string, error) {
	switch filename {
	case "Sigfaibles_effectif_siret.csv":
		return "effectif", nil
	case "Sigfaibles_debits.csv":
		return "debit", nil
	case "Sigfaibles_debits2.csv":
		return "debit", nil
	default:
		return "", errors.New("Unrecognized type for " + filename)
	}
}
