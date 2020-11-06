package prepareimport

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// CreateTempFiles creates a temporary directory with a batch of files, and clean up after the execution of tests
func CreateTempFiles(t *testing.T, batchkey BatchKey, filenames []string) string {
	t.Helper()
	parentDir, err := ioutil.TempDir(os.TempDir(), "example")
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Cleanup(func() { os.RemoveAll(parentDir) })

	batchDir := filepath.Join(parentDir, batchkey.String())
	os.Mkdir(batchDir, 0777)

	for _, filename := range filenames {
		tmpFilename := filepath.Join(batchDir, filename)
		if err := ioutil.WriteFile(tmpFilename, []byte{}, 0666); err != nil {
			t.Fatal(err.Error())
		}
	}

	return parentDir
}
