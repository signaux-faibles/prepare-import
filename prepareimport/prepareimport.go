package prepareimport

import (
	"errors"
	"fmt"
	"io"
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
	effectifFile, _ := filesProperty.GetEffectifFile()
	filterFile, _ := filesProperty.GetFilterFile()

	if effectifFile == nil && batchKey.IsSubBatch() {
		parentFilesProperty, _ := PopulateFilesProperty(pathname, newSafeBatchKey(batchKey.GetParentBatch()))
		effectifFile, _ = parentFilesProperty.GetEffectifFile()
	}

	if filterFile == nil && batchKey.IsSubBatch() {
		parentFilesProperty, _ := PopulateFilesProperty(pathname, newSafeBatchKey(batchKey.GetParentBatch()))
		filterFile, _ = parentFilesProperty.GetFilterFile()
	}

	// if needed, create a filter file from the effectif file
	if filterFile == nil {
		if effectifFile == nil {
			return nil, errors.New("filter is missing: batch should include a filter or one effectif file")
		}
		effectifFilePath := path.Join(pathname, effectifFile.Path())
		effectifBatch := effectifFile.BatchKey()
		filterFile = newBatchFile(effectifBatch, "filter_siren_"+effectifBatch.String()+".csv")
		if err = createFilterFromEffectif(path.Join(pathname, filterFile.Path()), effectifFilePath); err != nil {
			return nil, err
		}
		dateFinEffectif, err = createfilter.DetectDateFinEffectif(effectifFilePath, createfilter.DefaultNbIgnoredCols) // TODO: éviter de lire le fichier Effectif deux fois
		if err != nil {
			return nil, err
		}
	}

	// add the filter to filesProperty
	if filesProperty["filter"] == nil && filterFile != nil {
		if batchKey.IsSubBatch() {
			// copy the filter into the sub-batch's directory
			src := path.Join(pathname, filterFile.Path())
			dest := path.Join(pathname, batchKey.GetParentBatch(), batchKey.Path(), filterFile.Name())
			err = copy(src, dest)
			if err != nil {
				return nil, err
			}
			filterFile = newBatchFile(batchKey, filterFile.Name())
		}
		filesProperty["filter"] = append(filesProperty["filter"], filterFile)
	}

	// make sure we have date_fin_effectif
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

func createFilterFromEffectif(filterFilePath string, effectifFilePath string) error {
	if fileExists(filterFilePath) {
		return errors.New("about to overwrite existing filter file: " + filterFilePath)
	}
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

// Copy the src file to dst. Any existing file will be overwritten and will not
// copy file attributes. Source: https://stackoverflow.com/a/21061062/592254
func copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
