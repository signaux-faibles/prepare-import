package prepareimport

import (
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/signaux-faibles/prepare-import/createfilter"
)

// PrepareImport generates an Admin object from files found at given pathname of the file system.
func PrepareImport(pathname string, batchKey BatchKey, providedDateFinEffectif string) (AdminObject, error) {
	var err error
	filesProperty, unsupportedFiles := PopulateFilesProperty(pathname, batchKey)

	dateFinEffectif, err := checkOrCreateFilterFromEffectif(filesProperty, batchKey, pathname)
	if err != nil {
		return nil, err
	}

	if dateFinEffectif.IsZero() {
		dateFinEffectif, err = time.Parse("2006-01-02", providedDateFinEffectif)
		if err != nil {
			return nil, errors.New("date_fin_effectif is missing or invalid: " + providedDateFinEffectif)
		}
	}

	if len(unsupportedFiles) > 0 {
		err = UnsupportedFilesError{unsupportedFiles}
	}

	return AdminObject{
		"_id":            IDProperty{batchKey, "batch"},
		"files":          filesProperty,
		"complete_types": populateCompleteTypesProperty(filesProperty),
		"param":          populateParamProperty(batchKey, NewDateFinEffectif(dateFinEffectif)),
	}, err
}

func checkOrCreateFilterFromEffectif(filesProperty FilesProperty, batchKey BatchKey, pathname string) (dateFinEffectif time.Time, err error) {
	if filesProperty["filter"] == nil && filesProperty["effectif"] != nil {
		if len(filesProperty["effectif"]) != 1 {
			return dateFinEffectif, fmt.Errorf("generating a filter requires just 1 effectif file, found: %s", filesProperty["effectif"])
		}
		err = createAndAppendFilter(filesProperty, batchKey, pathname)
		if err != nil {
			return dateFinEffectif, err
		}
		effectifFile := path.Join(pathname, filesProperty["effectif"][0])
		// TODO: éviter de lire le fichier Effectif deux fois
		dateFinEffectif, err = createfilter.DetectDateFinEffectif(effectifFile, createfilter.DefaultNbIgnoredCols)
		if err != nil {
			return dateFinEffectif, err
		}
	}
	if filesProperty["filter"] == nil || len(filesProperty["filter"]) == 0 {
		err = errors.New("filter is missing: please include a filter or an effectif file")
	}
	return dateFinEffectif, err
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

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
