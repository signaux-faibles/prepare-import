package prepareimport

import (
	"errors"
	"io/ioutil"
	"os"
	"path"

	"github.com/signaux-faibles/prepare-import/createfilter"
)

// PrepareImport generates an Admin object from files found at given pathname of the file system.
func PrepareImport(pathname string, batchKey BatchKey, dateFinEffectif DateFinEffectif) (AdminObject, error) {
	batchPath := path.Join(pathname, batchKey.String())
	filenames, err := ReadFilenames(batchPath)
	if err != nil {
		return nil, err
	}
	augmentedFiles := []DataFile{}
	for _, file := range filenames {
		augmentedFiles = append(augmentedFiles, AugmentDataFile(file, batchPath))
	}
	adminObject, unsupportedFiles := PopulateAdminObject(augmentedFiles, batchKey, dateFinEffectif)

	filesProperty := adminObject["files"].(FilesProperty)
	if filesProperty["filter"] == nil && filesProperty["effectif"] != nil {
		err = createAndAppendFilter(filesProperty, batchKey, pathname)
		if err != nil {
			return nil, err
		}
	}
	if len(unsupportedFiles) > 0 {
		return adminObject, UnsupportedFilesError{unsupportedFiles}
	} else {
		return adminObject, nil
	}
}

func createAndAppendFilter(filesProperty FilesProperty, batchKey BatchKey, pathname string) error {
	filterFileName := path.Join(batchKey.Path(), "filter_siren_"+batchKey.String()+".csv")
	filterFilePath := path.Join(pathname, filterFileName)
	// make sure that file does not already exist
	if fileExists(filterFilePath) {
		return errors.New("warning: about to overwrite existing filter file: " + filterFilePath)
	}
	filesProperty["filter"] = append(filesProperty["filter"], filterFileName)
	filterWriter, err := os.Create(filterFilePath)
	if err != nil {
		return err
	}
	// TODO: make sure that there is only one effectif file
	return createfilter.CreateFilter(
		filterWriter,
		path.Join(pathname, filesProperty["effectif"][0]), // effectifFileName
		createfilter.DefaultNbMois,
		createfilter.DefaultMinEffectif,
		createfilter.DefaultNbIgnoredRecords,
	)
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

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
