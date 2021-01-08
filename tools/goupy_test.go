package tools

import (
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/signaux-faibles/prepare-import/prepareimport"
	"github.com/stretchr/testify/assert"
)

func TestGoupy(t *testing.T) {
	t.Run("should accept accentuated characters from info file", func(t *testing.T) {

		contentPerFile := map[string][]byte{
			"3bf7f0f805a66926452321b77ec0c1db": []byte{},
			"3bf7f0f805a66926452321b77ec0c1db.info": []byte(`{
				"ID":"3bf7f0f805a66926452321b77ec0c1db",
				"Size":12337702,
				"SizeIsDeferred":false,
				"Offset":0,
				"MetaData":{
					"filename":"Ellisphère-Tête de groupe-FinalV2-2015.xlsx",
					"filetype":"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
					"goup-path":"signauxfaibles",
					"private":"false"
				},
				"IsPartial":false,
				"IsFinal":false,
				"PartialUploads":null,
				"Storage":{
					"Path":"/var/lib/goup_base/tusd/3bf7f0f805a66926452321b77ec0c1db",
					"Type":"filestore"
				}
			}`),
		}
		batchKey, _ := prepareimport.NewBatchKey("1910")
		parentDir := prepareimport.CreateTempFilesWithContent(t, batchKey, contentPerFile)
		path := path.Join(parentDir, batchKey.Path())
		cmd := exec.Command("python", "goupy.py", path)
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		assert.NoError(t, err)
	})
}
