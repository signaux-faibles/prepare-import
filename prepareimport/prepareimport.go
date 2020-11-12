package prepareimport

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/signaux-faibles/prepare-import/createfilter"
)

// PrepareImport generates an Admin object from files found at given pathname of the file system.
func PrepareImport(pathname string, batchKey BatchKey, providedDateFinEffectif string) (AdminObject, error) {

	batchPath := getBatchPath(pathname, batchKey)
	if _, err := ioutil.ReadDir(path.Join(pathname, batchPath)); err != nil {
		return nil, fmt.Errorf("could not find directory %s in provided path", batchPath)
	}

	var err error
	filesProperty, unsupportedFiles := PopulateFilesProperty(pathname, batchKey)

	// To complete the FilesProperty, we need:
	// - a filter file (created from an effectif file, at the batch/parent level)
	// - a dateFinEffectif value (provided as parameter, or detected from effectif file)

	var dateFinEffectif time.Time
	if !filesProperty.HasFilterFile() {
		effectifFile, err := filesProperty.GetEffectifFile()
		if err != nil {
			return nil, errors.New("filter is missing: please include a filter or one effectif file")
		}
		filterFileName, err := createFilterFromEffectif(effectifFile, batchKey, pathname)
		if err != nil {
			return nil, err
		}
		filesProperty["filter"] = append(filesProperty["filter"], filterFileName)
		effectifFilePath := path.Join(pathname, filesProperty["effectif"][0])
		dateFinEffectif, err = createfilter.DetectDateFinEffectif(effectifFilePath, createfilter.DefaultNbIgnoredCols) // TODO: éviter de lire le fichier Effectif deux fois
		if err != nil {
			return nil, err
		}
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

func createFilterFromEffectif(effectifFile string, batchKey BatchKey, pathname string) (string, error) {
	filterFileName := path.Join(batchKey.Path(), "filter_siren_"+batchKey.String()+".csv")
	filterFilePath := path.Join(pathname, filterFileName)
	if fileExists(filterFilePath) {
		return "", errors.New("about to overwrite existing filter file: " + filterFilePath)
	}
	effectifFilePath := path.Join(pathname, effectifFile)
	err := createFilterFile(filterFilePath, effectifFilePath)
	return filterFileName, err
}

func createFilterFile(filterFilePath string, effectifFilePath string) error {
	filterWriter, err := os.Create(filterFilePath)
	if err != nil {
		return err
	}
	return createfilter.CreateFilter(
		filterWriter,     // output: the filter file
		effectifFilePath, // input: the effectif file
		createfilter.DefaultNbMois,
		createfilter.DefaultMinEffectif,
		createfilter.DefaultNbIgnoredCols,
	)
}

func getBatchPath(pathname string, batchKey BatchKey) string {
	if batchKey.IsSubBatch() {
		return path.Join(batchKey.GetParentBatch(), batchKey.String())
	}
	return batchKey.String()
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
