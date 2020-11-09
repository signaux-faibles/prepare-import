package prepareimport

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractFileTypeFromMetadata(t *testing.T) {

	t.Run("should return \"debit\" for bin file which original name included \"debits\"", func(t *testing.T) {
		got := ExtractFileTypeFromMetadata("9a047825d8173684b69994428449302f", MakeMetadata(MetadataProperty{
			"filename":  "Sigfaible_debits.csv",
			"goup-path": "urssaf",
		}))
		assert.Equal(t, debit, got)
	})

	t.Run("should return \"bdf\" for bin file which came from bdf", func(t *testing.T) {
		got := ExtractFileTypeFromMetadata("60d1bd320523904d8b8b427efbbd3928", MakeMetadata(MetadataProperty{
			"filename":  "FICHIER_SF_2020_02.csv",
			"goup-path": "bdf",
		}))
		assert.Equal(t, bdf, got)
	})

	t.Run("should return \"interim\" for bin file which had a sas7dbat extension", func(t *testing.T) {
		got := ExtractFileTypeFromMetadata("ab8613ab66ebddb2db21e36b92fc5b70", MakeMetadata(MetadataProperty{
			"filename":  "tab_19m10.sas7bdat",
			"goup-path": "dgefp",
		}))
		assert.Equal(t, interim, got)
	})
}