package prepareimport

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
	"time"

	"prepare-import/createfilter"
)

// PrepareImport generates an Admin object from files found at given pathname of the file system.
func PrepareImport(pathname string, batchKey BatchKey, providedDateFinEffectif string) (AdminObject, error) {

	batchPath := getBatchPath(pathname, batchKey)
	slog.Info("Liste les fichiers de données ", slog.Any("path", batchPath))
	if _, err := os.ReadDir(path.Join(pathname, batchPath)); err != nil {
		return AdminObject{}, fmt.Errorf("could not find directory %s in provided path", batchPath)
	}

	var err error
	filesProperty, unsupportedFiles := PopulateFilesProperty(pathname, batchKey)

	// To complete the FilesProperty, we need:
	// - a filter file (created from an effectif file, at the batch/parent level)
	// - a dateFinEffectif value (provided as parameter, or detected from effectif file)

	var dateFinEffectif time.Time
	effectifFile, _ := filesProperty.GetEffectifFile()
	filterFile, _ := filesProperty.GetFilterFile()
	sireneULFile, _ := filesProperty.GetSireneULFile()
	if (effectifFile == nil || filterFile == nil) && batchKey.IsSubBatch() {
		slog.Info("Looking for effectif and/or filter file", slog.String("batch", batchKey.GetParentBatch()))
		parentFilesProperty, _ := PopulateFilesProperty(pathname, newSafeBatchKey(batchKey.GetParentBatch()))
		if effectifFile == nil {
			effectifFile, _ = parentFilesProperty.GetEffectifFile()
		}
		if filterFile == nil {
			filterFile, _ = parentFilesProperty.GetFilterFile()
		}
	}

	if effectifFile != nil {
		slog.Info("Found effectif file ", slog.String("filename", effectifFile.Name()))
	}

	if filterFile != nil {
		slog.Info("Found filter file ", slog.String("filename", filterFile.Name()))
	}

	if sireneULFile != nil {
		slog.Info("Found sireneUL file", slog.String("filename", sireneULFile.Name()))
	}

	// if needed, create a filter file from the effectif file
	if filterFile == nil {
		if effectifFile == nil {
			return AdminObject{}, errors.New("filter is missing: batch should include a filter or one effectif file")
		}
		effectifFilePath := effectifFile.AbsolutePath(pathname)
		sireneULFilePath := sireneULFile.AbsolutePath(pathname)
		effectifBatch := effectifFile.BatchKey()
		filterFile = newBatchFile(effectifBatch, "filter_siren_"+effectifBatch.String()+".csv")
		slog.Info("Generating filter file", slog.String("filename", filterFile.Path()))
		if err = createFilterFromEffectifAndSirene(path.Join(pathname, filterFile.Path()), effectifFilePath, sireneULFilePath); err != nil {
			return AdminObject{}, err
		}
	}

	// add the filter to filesProperty
	if filesProperty["filter"] == nil && filterFile != nil {
		if batchKey.IsSubBatch() {
			// copy the filter into the sub-batch's directory
			slog.Info("Copying filter file", slog.String("filename", filterFile.Path()))
			src := path.Join(pathname, filterFile.Path())
			dest := path.Join(pathname, batchKey.GetParentBatch(), batchKey.Path(), filterFile.Name())
			err = copy(src, dest)
			if err != nil {
				return AdminObject{}, err
			}
			filterFile = newBatchFile(batchKey, filterFile.Name())
		}
		slog.Info("Adding filter file to batch", slog.String("filename", filterFile.Path()))
		filesProperty["filter"] = append(filesProperty["filter"], filterFile)
	}

	if effectifFile != nil {
		println("Detecting dateFinEffectif from effectif file ...")
		effectifFilePath := effectifFile.AbsolutePath(pathname)
		dateFinEffectif, err = createfilter.DetectDateFinEffectif(effectifFilePath, createfilter.DefaultNbIgnoredCols) // TODO: éviter de lire le fichier Effectif deux fois
		if err != nil {
			return AdminObject{}, err
		}
	}

	// make sure we have date_fin_effectif
	if dateFinEffectif.IsZero() {
		println("Still missing date_fin_effectif => parsing CLI parameter ...")
		dateFinEffectif, err = time.Parse("2006-01-02", providedDateFinEffectif)
		if err != nil {
			return AdminObject{}, errors.New("date_fin_effectif is missing or invalid: " + providedDateFinEffectif)
		}
	}

	if len(unsupportedFiles) > 0 {
		err = UnsupportedFilesError{unsupportedFiles}
	}

	return AdminObject{
		ID:            IDProperty{batchKey, "batch"},
		Files:         populateFilesPaths(filesProperty),
		CompleteTypes: populateCompleteTypesProperty(filesProperty),
		Param:         populateParamProperty(batchKey, NewDateFinEffectif(dateFinEffectif)),
	}, err
}

func createFilterFromEffectifAndSirene(filterFilePath string, effectifFilePath string, sireneULFilePath string) error {
	if fileExists(filterFilePath) {
		return errors.New("about to overwrite existing filter file: " + filterFilePath)
	}
	filterWriter, err := os.Create(filterFilePath)
	if err != nil {
		return err
	}
	categoriesJuridiqueFilter := createfilter.CategorieJuridiqueFilter(sireneULFilePath)

	return createfilter.CreateFilter(
		filterWriter,     // output: the filter file
		effectifFilePath, // input: the effectif file
		createfilter.DefaultNbMois,
		createfilter.DefaultMinEffectif,
		createfilter.DefaultNbIgnoredCols,
		categoriesJuridiqueFilter,
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
	return !os.IsNotExist(err)
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
