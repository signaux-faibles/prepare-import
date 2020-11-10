package prepareimport

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPopulateFilesProperty(t *testing.T) {
	t.Run("Should return an empty json when there is no file", func(t *testing.T) {
		filesProperty, unsupportedFiles := PopulateFilesPropertyFromDataFiles([]DataFile{}, dummyBatchKey.Path())
		assert.Len(t, unsupportedFiles, 0)
		assert.Equal(t, FilesProperty{}, filesProperty)
	})

	t.Run("PopulateFilesProperty should contain effectif file in \"effectif\" property", func(t *testing.T) {
		filesProperty, unsupportedFiles := PopulateFilesPropertyFromDataFiles([]DataFile{SimpleDataFile{"Sigfaibles_effectif_siret.csv"}}, dummyBatchKey.Path())
		if assert.Len(t, unsupportedFiles, 0) {
			assert.Equal(t, []string{dummyBatchKey.Path() + "Sigfaibles_effectif_siret.csv"}, filesProperty[effectif])
		}
	})

	t.Run("PopulateFilesProperty should contain one debit file in \"debit\" property", func(t *testing.T) {
		filesProperty, unsupportedFiles := PopulateFilesPropertyFromDataFiles([]DataFile{SimpleDataFile{"Sigfaibles_debits.csv"}}, dummyBatchKey.Path())
		expected := FilesProperty{debit: []string{dummyBatchKey.Path() + "Sigfaibles_debits.csv"}}
		assert.Len(t, unsupportedFiles, 0)
		assert.Equal(t, expected, filesProperty)
	})

	t.Run("PopulateFilesProperty should contain both debits files in \"debit\" property", func(t *testing.T) {
		filesProperty, unsupportedFiles := PopulateFilesPropertyFromDataFiles([]DataFile{SimpleDataFile{"Sigfaibles_debits.csv"}, SimpleDataFile{"Sigfaibles_debits2.csv"}}, dummyBatchKey.Path())
		if assert.Len(t, unsupportedFiles, 0) {
			assert.Equal(t, []string{dummyBatchKey.Path() + "Sigfaibles_debits.csv", dummyBatchKey.Path() + "Sigfaibles_debits2.csv"}, filesProperty[debit])
		}
	})

	t.Run("Should support multiple types of csv files", func(t *testing.T) {
		type File struct {
			Type     ValidFileType
			Filename string
		}
		files := []File{
			{"diane", "diane_req_2002.csv"},               // --> DIANE
			{"diane", "diane_req_dom_2002.csv"},           // --> DIANE
			{"effectif", "effectif_dom.csv"},              // --> EFFECTIF
			{"filter", "filter_siren_2002.csv"},           // --> FILTER
			{"sirene_ul", "sireneUL.csv"},                 // --> SIRENE_UL
			{"sirene", "StockEtablissement_utf8_geo.csv"}, // --> SIRENE
		}
		expectedFiles := FilesProperty{}
		inputFiles := []DataFile{}
		for _, file := range files {
			expectedFiles[file.Type] = append(expectedFiles[file.Type], dummyBatchKey.Path()+file.Filename)
			inputFiles = append(inputFiles, SimpleDataFile{file.Filename})
		}
		resFilesProperty, unsupportedFiles := PopulateFilesPropertyFromDataFiles(inputFiles, dummyBatchKey.Path())
		assert.Len(t, unsupportedFiles, 0)
		assert.Equal(t, expectedFiles, resFilesProperty)
	})

	t.Run("Should not include unsupported files", func(t *testing.T) {
		filesProperty, unsupportedFiles := PopulateFilesPropertyFromDataFiles([]DataFile{SimpleDataFile{"coco.csv"}}, dummyBatchKey.Path())
		assert.Len(t, unsupportedFiles, 1)
		assert.Equal(t, FilesProperty{}, filesProperty)
	})
	t.Run("Should report unsupported files", func(t *testing.T) {
		_, unsupportedFiles := PopulateFilesPropertyFromDataFiles([]DataFile{SimpleDataFile{"coco.csv"}}, dummyBatchKey.Path())
		assert.Equal(t, []string{dummyBatchKey.Path() + "coco.csv"}, unsupportedFiles)
	})
}
