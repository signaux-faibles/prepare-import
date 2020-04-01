package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"path/filepath"

	// "flag"
	// "fmt"
	"io/ioutil"
	"regexp"
	"strings"
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

type Filename interface {
	GetFilenameToImport() string // the name as it will be stored in Admin
	GetFiletype() string
}

type SimpleFilename struct {
	filename string
}

func (ffn SimpleFilename) GetFiletype() string {
	return GetFileType(ffn.filename)
}

func (ffn SimpleFilename) GetFilenameToImport() string {
	return ffn.filename
}

type UploadedFilename struct {
	filename string
	path     string
}

func (ffn UploadedFilename) GetFiletype() string {
	metaFilepath := filepath.Join(ffn.path, strings.Replace(ffn.filename, ".bin", ".info", 1))
	fileinfo, err := LoadMetadata(metaFilepath)
	if err != nil {
		log.Fatal(err)
	}
	filetype := GetFileTypeFromMetadata(metaFilepath, fileinfo)
	return filetype // e.g. "Sigfaible_debits.csv"
}

func (ffn UploadedFilename) GetFilenameToImport() string {
	return ffn.filename
}

func PrepareImport(pathname string) (AdminObject, error) {
	filenames, err := ReadFilenames(pathname)
	if err != nil {
		return nil, err
	}
	augmentedFiles := []Filename{}
	for _, file := range filenames {
		var filename Filename
		if strings.HasSuffix(file, ".bin") {
			filename = UploadedFilename{file, pathname}
		} else {
			filename = SimpleFilename{file}
		}
		augmentedFiles = append(augmentedFiles, filename)
	}
	return PurePrepareImport(augmentedFiles), nil
}

func PurePrepareImport(augmentedFilenames []Filename) AdminObject {
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

func PopulateFilesProperty(filenames []Filename) FilesProperty {
	filesProperty := FilesProperty{}
	for _, filename := range filenames {
		filetype := filename.GetFiletype()

		if filetype == "" {
			// Unsupported file
			continue
		}
		if _, exists := filesProperty[filetype]; !exists {
			filesProperty[filetype] = []string{}
		}
		filesProperty[filetype] = append(filesProperty[filetype], filename.GetFilenameToImport())
	}
	return filesProperty
}

var hasDianePrefix = regexp.MustCompile(`^diane`)
var mentionsEffectif = regexp.MustCompile(`effectif_`)
var mentionsDebits = regexp.MustCompile(`_debits`)
var hasFilterPrefix = regexp.MustCompile(`^filter_`)

type MetadataProperty map[string]string

type UploadedFileMeta struct {
	MetaData MetadataProperty
}

func GetFileTypeFromMetadata(filename string, fileinfo UploadedFileMeta) string {
	metadata := fileinfo.MetaData
	if metadata["goup-path"] == "bdf" {
		return "bdf"
	} else {
		return GetFileType(metadata["filename"])
	}
}

// GetFileType returns a file type from filename, or empty string for unsupported file names
func GetFileType(filename string) string {
	switch {
	case filename == "act_partielle_conso_depuis2014_FRANCE.csv":
		return "apconso"
	case filename == "act_partielle_ddes_depuis2015_FRANCE.csv":
		return "apdemande"
	case filename == "Sigfaible_etablissement_utf8.csv":
		return "admin_urssaf"
	case filename == "Sigfaible_effectif_siren.csv":
		return "effectif_ent"
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
	case strings.HasSuffix(filename, ".sas7bdat"):
		return "interim"
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
