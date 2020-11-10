package prepareimport

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPopulateAdminObject(t *testing.T) {
	t.Run("Should return the filename in the debit property", func(t *testing.T) {
		filename := SimpleDataFile{"Sigfaibles_debits.csv"}

		res, unsupported := PopulateAdminObject([]DataFile{filename}, dummyBatchKey, dummyDateFinEffectif)
		expected := FilesProperty{debit: []string{dummyBatchKey.Path() + "Sigfaibles_debits.csv"}}
		assert.Len(t, unsupported, 0)
		assert.Equal(t, expected, res["files"])
	})

	t.Run("Should return an empty complete_types property", func(t *testing.T) {
		filename := SimpleDataFile{"Sigfaibles_debits.csv"}

		res, unsupported := PopulateAdminObject([]DataFile{filename}, dummyBatchKey, dummyDateFinEffectif)
		expected := []ValidFileType{}
		assert.Len(t, unsupported, 0)
		assert.Equal(t, expected, res["complete_types"])
	})

	t.Run("Should return apconso as a complete_type", func(t *testing.T) {
		filename := SimpleDataFile{"act_partielle_conso_depuis2014_FRANCE.csv"}
		res, unsupported := PopulateAdminObject([]DataFile{filename}, dummyBatchKey, dummyDateFinEffectif)
		expected := []ValidFileType{apconso}
		assert.Len(t, unsupported, 0)
		assert.Equal(t, expected, res["complete_types"])
	})

	t.Run("Should return an empty json when there is no file", func(t *testing.T) {
		res, unsupported := PopulateAdminObject([]DataFile{}, dummyBatchKey, dummyDateFinEffectif)
		assert.Len(t, unsupported, 0)
		assert.Equal(t, FilesProperty{}, res["files"])
	})

	t.Run("Should support multiple types of csv files", func(t *testing.T) {
		files := []string{
			"diane_req_2002.csv",              // --> DIANE
			"diane_req_dom_2002.csv",          // --> DIANE
			"effectif_dom.csv",                // --> EFFECTIF
			"filter_siren_2002.csv",           // --> FILTER
			"sireneUL.csv",                    // --> SIRENE_UL
			"StockEtablissement_utf8_geo.csv", // --> SIRENE
		}
		expectedFiles := []string{}
		augmentedFiles := []DataFile{}
		for _, file := range files {
			expectedFiles = append(expectedFiles, dummyBatchKey.Path()+file)
			augmentedFiles = append(augmentedFiles, SimpleDataFile{file})
		}
		res, unsupported := PopulateAdminObject(augmentedFiles, dummyBatchKey, dummyDateFinEffectif)
		assert.Len(t, unsupported, 0)
		resFilesProperty := res["files"].(FilesProperty)
		resultingFiles := []string{}
		for _, filenames := range resFilesProperty {
			resultingFiles = append(resultingFiles, filenames...)
		}
		assert.Subset(t, resultingFiles, expectedFiles)
	})

	t.Run("Should return an _id property", func(t *testing.T) {
		res, unsupported := PopulateAdminObject([]DataFile{}, newSafeBatchKey("1802"), dummyDateFinEffectif)
		assert.Len(t, unsupported, 0)
		assert.Equal(t, IDProperty{newSafeBatchKey("1802"), "batch"}, res["_id"])
	})

	t.Run("Should return a date_fin consistent with batch key", func(t *testing.T) {
		res, unsupported := PopulateAdminObject([]DataFile{}, newSafeBatchKey("1912"), dummyDateFinEffectif) // ~= 12/2019
		expected := MongoDate{"2019-12-01T00:00:00.000+0000"}
		assert.Len(t, unsupported, 0)
		assert.Equal(t, expected, res["param"].(ParamProperty).DateFin)
	})
}
