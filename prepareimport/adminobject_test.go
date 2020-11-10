package prepareimport

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPopulateAdminObject(t *testing.T) {
	t.Run("Should return the filename in the debit property", func(t *testing.T) {
		filename := SimpleDataFile{"Sigfaibles_debits.csv"}

		res, unsupported := PopulateAdminObject([]DataFile{filename}, dummyBatchKey, validDateFinEffectif)
		expected := FilesProperty{debit: []string{dummyBatchKey.Path() + "Sigfaibles_debits.csv"}}
		assert.Len(t, unsupported, 0)
		assert.Equal(t, expected, res["files"])
	})

	t.Run("Should return an empty complete_types property", func(t *testing.T) {
		filename := SimpleDataFile{"Sigfaibles_debits.csv"}

		res, unsupported := PopulateAdminObject([]DataFile{filename}, dummyBatchKey, validDateFinEffectif)
		expected := []ValidFileType{}
		assert.Len(t, unsupported, 0)
		assert.Equal(t, expected, res["complete_types"])
	})

	t.Run("Should return apconso as a complete_type", func(t *testing.T) {
		filename := SimpleDataFile{"act_partielle_conso_depuis2014_FRANCE.csv"}
		res, unsupported := PopulateAdminObject([]DataFile{filename}, dummyBatchKey, validDateFinEffectif)
		expected := []ValidFileType{apconso}
		assert.Len(t, unsupported, 0)
		assert.Equal(t, expected, res["complete_types"])
	})

	t.Run("Should return an _id property", func(t *testing.T) {
		res, unsupported := PopulateAdminObject([]DataFile{}, newSafeBatchKey("1802"), validDateFinEffectif)
		assert.Len(t, unsupported, 0)
		assert.Equal(t, IDProperty{newSafeBatchKey("1802"), "batch"}, res["_id"])
	})
}

func TestPopulateParamProperty(t *testing.T) {
	t.Run("Should return a date_fin consistent with batch key", func(t *testing.T) {
		res := populateParamProperty(newSafeBatchKey("1912"), validDateFinEffectif)
		expected := MongoDate{"2019-12-01T00:00:00.000+0000"}
		assert.Equal(t, expected, res.DateFin)
	})
}
