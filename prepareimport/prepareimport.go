package prepareimport

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/signaux-faibles/prepare-import/createfilter"
)

// PrepareImport generates an Admin object from files found at given pathname of the file system.
func PrepareImport(pathname string, batchKey BatchKey, providedDateFinEffectif string) (AdminObject, error) {

	batchPath := path.Join(pathname, batchKey.String())
	filenames, err := ReadFilenames(batchPath)
	if err != nil {
		return nil, err
	}
	augmentedFiles := []DataFile{}
	for _, file := range filenames {
		augmentedFiles = append(augmentedFiles, AugmentDataFile(file, batchPath))
	}

	// adminObject, unsupportedFiles := PopulateAdminObject(augmentedFiles, batchKey, dateFinEffectif) // TODO: de-duplicate code
	filesProperty, unsupportedFiles := PopulateFilesProperty(augmentedFiles, batchKey.Path())
	var completeTypes = []ValidFileType{}
	for _, typeName := range defaultCompleteTypes {
		if _, ok := filesProperty[typeName]; ok {
			completeTypes = append(completeTypes, typeName)
		}
	}

	var dateFinEffectif time.Time
	if filesProperty["filter"] == nil && filesProperty["effectif"] != nil {
		err = createAndAppendFilter(filesProperty, batchKey, pathname)
		if err != nil {
			return nil, err
		}
		dateFinEffectif = time.Date(2020, time.Month(1), 1, 0, 0, 0, 0, time.UTC) // TODO: detect from file
	}

	if filesProperty["filter"] == nil || len(filesProperty["filter"]) == 0 {
		return nil, errors.New("filter is missing: please include a filter or an effectif file")
	}

	if dateFinEffectif.IsZero() {
		dateFinEffectif, err = time.Parse("2006-01-02", providedDateFinEffectif)
		if err != nil {
			return nil, errors.New("date_fin_effectif is missing or invalid: " + providedDateFinEffectif)
		}
	}

	adminObject := AdminObject{
		"_id":            IDProperty{batchKey, "batch"},
		"files":          filesProperty,
		"complete_types": completeTypes,
		"param": ParamProperty{
			DateDebut:       MongoDate{"2014-01-01T00:00:00.000+0000"},
			DateFin:         MongoDate{"20" + batchKey.String()[0:2] + "-" + batchKey.String()[2:4] + "-01T00:00:00.000+0000"},
			DateFinEffectif: NewDateFinEffectif(dateFinEffectif).MongoDate(),
		},
	}

	if len(unsupportedFiles) > 0 {
		return adminObject, UnsupportedFilesError{unsupportedFiles}
	}
	return adminObject, nil
}

func createAndAppendFilter(filesProperty FilesProperty, batchKey BatchKey, pathname string) error {
	// make sure that there is only one effectif file
	if len(filesProperty["effectif"]) != 1 {
		return errors.New("filter generation requires just 1 effectif file")
	}
	// create the filter file, if it does not already exist
	filterFileName := path.Join(batchKey.Path(), "filter_siren_"+batchKey.String()+".csv")
	filterFilePath := path.Join(pathname, filterFileName)
	if fileExists(filterFilePath) {
		return errors.New("about to overwrite existing filter file: " + filterFilePath)
	}
	filterWriter, err := os.Create(filterFilePath)
	if err != nil {
		return err
	}
	// populate the filter file and the "filter" property of the AdminObject
	err = createfilter.CreateFilter(
		filterWriter,
		path.Join(pathname, filesProperty["effectif"][0]), // effectifFileName
		createfilter.DefaultNbMois,
		createfilter.DefaultMinEffectif,
		createfilter.DefaultNbIgnoredCols,
	)
	if err != nil {
		return err
	}
	filesProperty["filter"] = append(filesProperty["filter"], filterFileName)
	return nil
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
