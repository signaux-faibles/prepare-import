package prepareimport

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPopulateFilesProperty(t *testing.T) {
	t.Run("PopulateFilesProperty should contain effectif file in \"effectif\" property", func(t *testing.T) {
		filesProperty, unsupportedFiles := PopulateFilesProperty([]DataFile{SimpleDataFile{"Sigfaibles_effectif_siret.csv"}}, DUMMY_BATCHKEY.Path())
		if assert.Len(t, unsupportedFiles, 0) {
			assert.Equal(t, []string{DUMMY_BATCHKEY.Path() + "Sigfaibles_effectif_siret.csv"}, filesProperty[EFFECTIF])
		}
	})

	t.Run("PopulateFilesProperty should contain one debit file in \"debit\" property", func(t *testing.T) {
		filesProperty, unsupportedFiles := PopulateFilesProperty([]DataFile{SimpleDataFile{"Sigfaibles_debits.csv"}}, DUMMY_BATCHKEY.Path())
		if assert.Len(t, unsupportedFiles, 0) {
			assert.Equal(t, []string{DUMMY_BATCHKEY.Path() + "Sigfaibles_debits.csv"}, filesProperty[DEBIT])
		}
	})

	t.Run("PopulateFilesProperty should contain both debits files in \"debit\" property", func(t *testing.T) {
		filesProperty, unsupportedFiles := PopulateFilesProperty([]DataFile{SimpleDataFile{"Sigfaibles_debits.csv"}, SimpleDataFile{"Sigfaibles_debits2.csv"}}, DUMMY_BATCHKEY.Path())
		if assert.Len(t, unsupportedFiles, 0) {
			assert.Equal(t, []string{DUMMY_BATCHKEY.Path() + "Sigfaibles_debits.csv", DUMMY_BATCHKEY.Path() + "Sigfaibles_debits2.csv"}, filesProperty[DEBIT])
		}
	})

	t.Run("Should not include unsupported files", func(t *testing.T) {
		filesProperty, unsupportedFiles := PopulateFilesProperty([]DataFile{SimpleDataFile{"coco.csv"}}, DUMMY_BATCHKEY.Path())
		assert.Len(t, unsupportedFiles, 1)
		assert.Equal(t, FilesProperty{}, filesProperty)
	})
	t.Run("Should report unsupported files", func(t *testing.T) {
		_, unsupportedFiles := PopulateFilesProperty([]DataFile{SimpleDataFile{"coco.csv"}}, DUMMY_BATCHKEY.Path())
		assert.Equal(t, []string{DUMMY_BATCHKEY.Path() + "coco.csv"}, unsupportedFiles)
	})
}
