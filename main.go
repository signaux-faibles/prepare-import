package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
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
	// TODO: serialize to JSON string
}

type AdminObject map[string]interface{}

func PrepareImport(pathname string) (AdminObject, error) {
	filenames, err := ReadFilenames(pathname)
	if err != nil {
		return nil, err
	}
	return PurePrepareImport(filenames), nil
}

func PurePrepareImport(filenames []string) AdminObject {
	filesProperty := PopulateFilesProperty(filenames)
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

func PopulateFilesProperty(filenames []string) FilesProperty {
	filesProperty := FilesProperty{
		// "effectif": []string{"coucou"},
		// "debit":    []string{},
	}
	for _, filename := range filenames {
		filetype := GetFileType(filename)
		if filetype == "" {
			// Unsupported file
			continue
		}
		if _, exists := filesProperty[filetype]; !exists {
			filesProperty[filetype] = []string{}
		}
		filesProperty[filetype] = append(filesProperty[filetype], filename)
	}
	return filesProperty
}

var hasDianePrefix = regexp.MustCompile(`^diane`)
var mentionsEffectif = regexp.MustCompile(`effectif_`)
var hasFilterPrefix = regexp.MustCompile(`^filter_`)

// GetFileType returns a file type from filename, or empty string for unsupported file names
func GetFileType(filename string) string {
	switch {
	case hasDianePrefix.MatchString(filename):
		return "diane"
	case mentionsEffectif.MatchString(filename):
		return "effectif"
	case hasFilterPrefix.MatchString(filename):
		return "filter"
	case filename == "sireneUL.csv":
		return "sirene_ul"
	case filename == "StockEtablissement_utf8_geo.csv":
		return "comptes"
	case filename == "Sigfaibles_debits.csv":
		return "debit"
	case filename == "Sigfaibles_debits2.csv":
		return "debit"
	default:
		return ""
	}
}
