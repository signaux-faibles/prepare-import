package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
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

func DefaultMetadataReader(filename string) UploadedFileMeta {
	return UploadedFileMeta{} // TODO
}

func PopulateFilesProperty(filenames []string) FilesProperty {
	filesProperty := FilesProperty{
		// "effectif": []string{"coucou"},
		// "debit":    []string{},
	}
	for _, filename := range filenames {
		filetype := GetFileType(filename, DefaultMetadataReader)
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
var mentionsDebits = regexp.MustCompile(`_debits`)
var hasFilterPrefix = regexp.MustCompile(`^filter_`)

type UploadedFileMeta map[string]interface{}

// GetFileType returns a file type from filename, or empty string for unsupported file names
func GetFileType(filename string, getFileMeta func(string) UploadedFileMeta) string {
	switch {
	case strings.HasSuffix(filename, ".bin"):
		metadata := getFileMeta(filename)["MetaData"].(map[string]string)
		if metadata["goup-path"] == "bdf" {
			return "bdf"
		}
		return GetFileType(metadata["filename"], DefaultMetadataReader)
	case filename == "Sigfaible_pcoll.csv":
		return "procol"
	case filename == "Sigfaible_cotisdues.csv":
		return "cotisation"
	case filename == "Sigfaible_delais.csv":
		return "delai"
	case filename == "Sigfaible_ccsf.csv":
		return "ccsf"
	case filename == "sireneUL.csv":
		return "sirene_ul"
	case filename == "StockEtablissement_utf8_geo.csv":
		return "comptes"
	case mentionsDebits.MatchString(filename):
		return "debit"
	case hasDianePrefix.MatchString(filename):
		return "diane"
	case mentionsEffectif.MatchString(filename):
		return "effectif"
	case hasFilterPrefix.MatchString(filename):
		return "filter"
	default:
		return ""
	}
}
