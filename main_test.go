// $ go test # to run the tests

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestMain(t *testing.T) {
// 	t.Run("Should return a files property from a directory name", func(t *testing.T) {
// 		res := main()
// 		expected := "{\"files\": {\"debit\": [\"Sigfaibles_debits.csv\"]}}"
// 		t.Equal(t, expected, res)
// 	})
// }

// Test: ReadFilenames should return filenames in a directory
func TestReadFilenames(t *testing.T) {
	t.Run("Should return filenames in a directory", func(t *testing.T) {
		dir, err := ioutil.TempDir(os.TempDir(), "example")
		if err != nil {
			t.Fatal(err.Error())
		}

		defer os.RemoveAll(dir) // clean up

		tmpFilename := filepath.Join(dir, "tmpfile")
		if err := ioutil.WriteFile(tmpFilename, []byte{}, 0666); err != nil {
			t.Fatal(err.Error())
		}

		filenames, err := ReadFilenames(dir)
		if err != nil {
			t.Fatal(err.Error())
		}
		assert.Equal(t, []string{"tmpfile"}, filenames)
	})
}

// Prepare import should return json object.
func TestPrepareImport(t *testing.T) {
	t.Run("Should return a json with one file", func(t *testing.T) {
		res, _ := PrepareImport([]string{"Sigfaibles_debits.csv"})
    expected := map[string]interface{}{
      "files":FileProperty{"debit": []string{"Sigfaibles_debits.csv"}},
    }
		assert.Equal(t, expected, res)
	})
	t.Run("Should return an empty json when there is no file", func(t *testing.T) {
		res, _ := PrepareImport([]string{})
    assert.Equal(t, map[string]interface{}{"files": FileProperty{}}, res)
	})
}

func TestPopulateFilesProperty(t *testing.T) {

	// t.Run() is used to define sub-tests. (see https://golang.org/pkg/testing/#hdr-Subtests_and_Sub_benchmarks)

	t.Run("PopulateFilesProperty should contain effectif file in \"effectif\" property", func(t *testing.T) {
		filesProperty := PopulateFilesProperty([]string{"Sigfaibles_effectif_siret.csv"})
		// isEqualSlice(t, filesProperty["effectif"], []string{"Sigfaibles_effectif_siret.csv"})
		assert.Equal(t, []string{"Sigfaibles_effectif_siret.csv"}, filesProperty["effectif"])
	})

	t.Run("PopulateFilesProperty should contain one debit file in \"debit\" property", func(t *testing.T) {
		filesProperty := PopulateFilesProperty([]string{"Sigfaibles_debits.csv"})
		assert.Equal(t, []string{"Sigfaibles_debits.csv"}, filesProperty["debit"])
	})

	t.Run("PopulateFilesProperty should contain both debits files in \"debit\" property", func(t *testing.T) {
		filesProperty := PopulateFilesProperty([]string{"Sigfaibles_debits.csv", "Sigfaibles_debits2.csv"})
		assert.Equal(t, []string{"Sigfaibles_debits.csv", "Sigfaibles_debits2.csv"}, filesProperty["debit"])
	})
}

func TestGetFileType(t *testing.T) {
	t.Run("should return \"effectif\" for \"Sigfaibles_effectif_siret.csv\"", func(t *testing.T) {
		got, _ := GetFileType("Sigfaibles_effectif_siret.csv")
		assert.Equal(t, "effectif", got)
	})

	t.Run("should return \"debit\" for \"Sigfaibles_debits.csv\"", func(t *testing.T) {
		got, _ := GetFileType("Sigfaibles_debits.csv")
		assert.Equal(t, "debit", got)
	})

	t.Run("should return \"debit\" for \"Sigfaibles_debits2.csv\"", func(t *testing.T) {
		got, _ := GetFileType("Sigfaibles_debits2.csv")
		assert.Equal(t, "debit", got)
	})
}
