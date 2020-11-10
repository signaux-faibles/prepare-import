package prepareimport

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func newSafeBatchKey(key string) BatchKey {
	batchKey, err := NewBatchKey(key)
	if err != nil {
		log.Fatal(err)
	}
	return batchKey
}

var dummyBatchKey = newSafeBatchKey("1802")

var dummyDateFinEffectif = "2014-01-01"
var validDateFinEffectif = NewDateFinEffectif(time.Date(2014, time.January, 1, 0, 0, 0, 0, time.UTC)) // "2014-01-01"

// CreateTempFiles creates a temporary directory with a batch of files, and clean up after the execution of tests
func CreateTempFiles(t *testing.T, batchkey BatchKey, filenames []string) string {
	contentPerFile := map[string][]byte{}
	for _, filename := range filenames {
		contentPerFile[filename] = []byte{}
	}
	return CreateTempFilesWithContent(t, batchkey, contentPerFile)
}

// CreateTempFilesWithContent creates a temporary directory with a batch of files, and clean up after the execution of tests
func CreateTempFilesWithContent(t *testing.T, batchkey BatchKey, contentPerFile map[string][]byte) string {
	t.Helper()
	parentDir, err := ioutil.TempDir(os.TempDir(), "example")
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Cleanup(func() { os.RemoveAll(parentDir) })

	batchDir := filepath.Join(parentDir, batchkey.String())
	os.Mkdir(batchDir, 0777)

	for filename := range contentPerFile {
		tmpFilename := filepath.Join(batchDir, filename)
		if err := ioutil.WriteFile(tmpFilename, contentPerFile[filename], 0666); err != nil {
			t.Fatal(err.Error())
		}
	}

	return parentDir
}
