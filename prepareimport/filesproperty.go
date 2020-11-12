package prepareimport

import (
	"io/ioutil"
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
