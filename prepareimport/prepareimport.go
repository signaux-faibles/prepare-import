package prepareimport

import (
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
	adminObject, err := PopulateAdminObject(augmentedFiles, batchKey, dateFinEffectif)
	if err != nil {
		return nil, err
	}
	filesProperty := adminObject["files"].(FilesProperty)
	if filesProperty["filter"] == nil && filesProperty["effectif"] != nil {
		filterFileName := path.Join(batchKey.Path(), "filter_siren_"+batchKey.String()+".csv")
		filesProperty["filter"] = append(filesProperty["filter"], filterFileName)
		// TODO: make sure that file does not already exist
		filterWriter, err := os.Create(path.Join(pathname, filterFileName))
		if err != nil {
			return nil, err
		}
		// TODO: make sure that there is only one effectif file
		err = createfilter.CreateFilter(
			filterWriter,
			path.Join(pathname, filesProperty["effectif"][0]), // effectifFileName
			createfilter.DefaultNbMois,
			createfilter.DefaultMinEffectif,
			createfilter.DefaultNbIgnoredRecords,
		)
		if err != nil {
			return nil, err
		}
	}
	return adminObject, nil
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
