package prepareimport

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func dummyBatchFile(filename string) BatchFile {
	return newBatchFile(dummyBatchKey, filename)
}

func TestPopulateCompleteTypesProperty(t *testing.T) {
	t.Run("Should not return a small debit file as a complete_type", func(t *testing.T) {
		res := populateCompleteTypesProperty(FilesProperty{"debit": {dummyBatchFile("Sigfaibles_debits.csv")}})
		expected := []ValidFileType{}
		assert.Equal(t, expected, res)
	})

	t.Run("Should return a large gzipped debit file as a complete_type", func(t *testing.T) {
		expected := []ValidFileType{debit}
		debitBatchFile := batchFile{
			batchKey:    dummyBatchKey,
			filename:    "Sigfaibles_debits.csv",
			gzippedSize: 254781489, // thresholdPerGzippedFileType["debit"]
		}
		res := populateCompleteTypesProperty(FilesProperty{"debit": {debitBatchFile}})
		assert.Equal(t, expected, res)
	})

	t.Run("Should return apconso as a complete_type", func(t *testing.T) {
		res := populateCompleteTypesProperty(FilesProperty{"apconso": {dummyBatchFile("act_partielle_conso_depuis2014_FRANCE.csv")}})
		expected := []ValidFileType{apconso}
		assert.Equal(t, expected, res)
	})

}

func TestPopulateParamProperty(t *testing.T) {
	t.Run("Should return a date_fin consistent with batch key", func(t *testing.T) {
		res := populateParamProperty(newSafeBatchKey("1912"), validDateFinEffectif)
		expected := MongoDate{"2019-12-01T00:00:00.000+0000"}
		assert.Equal(t, expected, res.DateFin)
	})
}
