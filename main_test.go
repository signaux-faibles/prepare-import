package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTempFiles(t *testing.T, filename string) string {
	t.Helper()
	dir, err := ioutil.TempDir(os.TempDir(), "example")
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Cleanup(func() { os.RemoveAll(dir) })

	tmpFilename := filepath.Join(dir, filename)
	if err := ioutil.WriteFile(tmpFilename, []byte{}, 0666); err != nil {
		t.Fatal(err.Error())
	}

	return dir
}

// Test: ReadFilenames should return filenames in a directory
func TestReadFilenames(t *testing.T) {
	t.Run("Should return filenames in a directory", func(t *testing.T) {
		dir := createTempFiles(t, "tmpfile")
		filenames, err := ReadFilenames(dir)
		if err != nil {
			t.Fatal(err.Error())
		}
		assert.Equal(t, []string{"tmpfile"}, filenames)
	})
}

func TestPrepareImport(t *testing.T) {
	t.Run("Should return a json with one file", func(t *testing.T) {
		dir := createTempFiles(t, "Sigfaibles_debits.csv")
		res, _ := PrepareImport(dir)
		expected := AdminObject{
			"files": FilesProperty{"debit": []string{"Sigfaibles_debits.csv"}},
		}
		assert.Equal(t, expected, res)
	})
}

// Prepare import should return json object.
func TestPurePrepareImport(t *testing.T) {
	t.Run("Should return a json with one file", func(t *testing.T) {
		res := PurePrepareImport([]string{"Sigfaibles_debits.csv"})
		expected := AdminObject{
			"files": FilesProperty{"debit": []string{"Sigfaibles_debits.csv"}},
		}
		assert.Equal(t, expected, res)
	})

	t.Run("Should return an empty json when there is no file", func(t *testing.T) {
		res := PurePrepareImport([]string{})
		assert.Equal(t, AdminObject{"files": FilesProperty{}}, res)
	})

	t.Run("Should support multiple types of csv files", func(t *testing.T) {
		files := []string{
			"diane_req_2002.csv",              // --> "diane"
			"diane_req_dom_2002.csv",          // --> "diane"
			"effectif_dom.csv",                // --> "effectif"
			"filter_siren_2002.csv",           // --> "filter"
			"sireneUL.csv",                    // --> "sirene_ul"
			"StockEtablissement_utf8_geo.csv", // --> "comptes"
		}
		res := PurePrepareImport(files)
		resFilesProperty := res["files"].(FilesProperty)
		resultingFiles := []string{}
		for _, filenames := range resFilesProperty {
			resultingFiles = append(resultingFiles, filenames...)
		}
		assert.Subset(t, resultingFiles, files)
	})
}

func TestPopulateFilesProperty(t *testing.T) {
	t.Run("PopulateFilesProperty should contain effectif file in \"effectif\" property", func(t *testing.T) {
		filesProperty := PopulateFilesProperty([]string{"Sigfaibles_effectif_siret.csv"})
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

	t.Run("Should not include unsupported files", func(t *testing.T) {
		filesProperty := PopulateFilesProperty([]string{"coco.csv"})
		assert.Equal(t, FilesProperty{}, filesProperty)
	})
}

func TestGetFileType(t *testing.T) {
	t.Run("should return \"effectif\" for \"Sigfaibles_effectif_siret.csv\"", func(t *testing.T) {
		got := GetFileType("Sigfaibles_effectif_siret.csv", DefaultMetadataReader)
		assert.Equal(t, "effectif", got)
	})

	t.Run("should return \"debit\" for \"Sigfaibles_debits.csv\"", func(t *testing.T) {
		got := GetFileType("Sigfaibles_debits.csv", DefaultMetadataReader)
		assert.Equal(t, "debit", got)
	})

	t.Run("should return \"debit\" for \"Sigfaibles_debits2.csv\"", func(t *testing.T) {
		got := GetFileType("Sigfaibles_debits2.csv", DefaultMetadataReader)
		assert.Equal(t, "debit", got)
	})

	t.Run("should return \"urssaf\" for bin file which come from urssaf", func(t *testing.T) {
		got := GetFileType("15b6ceeb928a3bc160b0e2dc2a794ad4.bin", func(filename string) UploadedFileMeta {
			return UploadedFileMeta{
				"MetaData": map[string]string{
					"filename":  "Sigfaible_cotisdues.csv",
					"filetype":  "application/vnd.ms-excel",
					"goup-path": "urssaf",
					"private":   "false",
				},
			}
		})
		assert.Equal(t, "urssaf", got)
	})

	// inspired by https://github.com/golang/go/wiki/TableDrivenTests
	cases := []struct {
		name     string
		category string
	}{
		{"diane_req_2002.csv", "diane"},
		{"diane_req_dom_2002.csv", "diane"},
		{"effectif_dom.csv", "effectif"},
		{"filter_siren_2002.csv", "filter"},
		{"sireneUL.csv", "sirene_ul"},
		{"StockEtablissement_utf8_geo.csv", "comptes"},
	}
	for _, testCase := range cases {
		t.Run("should return "+testCase.category+" for file "+testCase.name, func(t *testing.T) {
			got := GetFileType(testCase.name, DefaultMetadataReader)
			assert.Equal(t, testCase.category, got)
		})
	}
}
