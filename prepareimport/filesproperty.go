package prepareimport

import (
	"path"
	"strings"
)

// FilesProperty represents the "files" property of an Admin object.
type FilesProperty map[ValidFileType][]string

// PopulateFilesProperty populates the "files" property of an Admin object, given a path.
func PopulateFilesProperty(pathname string, batchKey BatchKey) (FilesProperty, []string) {
	batchPath := path.Join(pathname, batchKey.String())
	filenames, _ := ReadFilenames(batchPath)
	augmentedFiles := []DataFile{}
	for _, file := range filenames {
		augmentedFiles = append(augmentedFiles, AugmentDataFile(file, batchPath))
	}

	return PopulateFilesPropertyFromDataFiles(augmentedFiles, batchKey.Path())
}

// PopulateFilesPropertyFromDataFiles populates the "files" property of an Admin object, given a list of Data files.
func PopulateFilesPropertyFromDataFiles(filenames []DataFile, prefix string) (FilesProperty, []string) {
	filesProperty := FilesProperty{}
	unsupportedFiles := []string{}
	for _, filename := range filenames {
		filetype := filename.DetectFileType()

		if filetype == "" {
			if !strings.HasSuffix(filename.GetFilename(), ".info") {
				unsupportedFiles = append(unsupportedFiles, prefix+filename.GetFilename())
			}
			continue
		}
		if _, exists := filesProperty[filetype]; !exists {
			filesProperty[filetype] = []string{}
		}
		filesProperty[filetype] = append(filesProperty[filetype], prefix+filename.GetFilename())
	}
	return filesProperty, unsupportedFiles
}

// ParamProperty represents the "param" property of an Admin object.
type ParamProperty struct {
	DateDebut       MongoDate `json:"date_debut"`
	DateFin         MongoDate `json:"date_fin"`
	DateFinEffectif MongoDate `json:"date_fin_effectif"`
}
