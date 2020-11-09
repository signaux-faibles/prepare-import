package prepareimport

import (
	"io/ioutil"
	"path"
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
	if filesProperty["effectif"] != nil && filesProperty["effectif_ent"] != nil {
		filterFile := path.Join(batchKey.Path(), "filter_siren_"+batchKey.String()+".csv")
		filesProperty["filter"] = append(filesProperty["filter"], filterFile)
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
