package prepareimport

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPopulateCompleteTypesProperty(t *testing.T) {
	t.Run("Should not return a debit file as a complete_type", func(t *testing.T) {
		res := populateCompleteTypesProperty(FilesProperty{"debit": {"Sigfaibles_debits.csv"}})
		expected := []ValidFileType{}
		assert.Equal(t, expected, res)
	})

	t.Run("Should return apconso as a complete_type", func(t *testing.T) {
		res := populateCompleteTypesProperty(FilesProperty{"apconso": {"act_partielle_conso_depuis2014_FRANCE.csv"}})
		expected := []ValidFileType{apconso}
		assert.Equal(t, expected, res)
	})

}

func TestPopulateAdminObject(t *testing.T) {

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
