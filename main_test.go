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

func MakeMetadataReader(metadataFields map[string]string) func(string) UploadedFileMeta {
	return func(filename string) UploadedFileMeta {
		return UploadedFileMeta{"MetaData": metadataFields}
	}
}

func DummyMetadataReader(filename string) UploadedFileMeta {
	return UploadedFileMeta{}
}

func TestGetFileType(t *testing.T) {

	t.Run("should return \"debit\" for bin file which original name included \"debits\"", func(t *testing.T) {
		got := GetFileType("9a047825d8173684b69994428449302f.bin", MakeMetadataReader(map[string]string{
			"filename":  "Sigfaible_debits.csv",
			"goup-path": "urssaf",
		}))
		assert.Equal(t, "debit", got)
	})

	t.Run("should return \"bdf\" for bin file which came from bdf", func(t *testing.T) {
		got := GetFileType("60d1bd320523904d8b8b427efbbd3928.bin", MakeMetadataReader(map[string]string{
			"filename":  "FICHIER_SF_2020_02.csv",
			"goup-path": "bdf",
		}))
		assert.Equal(t, "bdf", got)
	})

	t.Run("should return \"interim\" for bin file which had a sas7dbat extension", func(t *testing.T) {
		got := GetFileType("ab8613ab66ebddb2db21e36b92fc5b70.bin", MakeMetadataReader(map[string]string{
			"filename":  "tab_19m10.sas7bdat",
			"goup-path": "dgefp",
		}))
		assert.Equal(t, "interim", got)
	})

	// inspired by https://github.com/golang/go/wiki/TableDrivenTests
	cases := []struct {
		name     string
		category string
	}{
		// guessed from urssaf files found on stockage/goub server
		{"Sigfaible_debits.csv", "debit"},
		{"Sigfaible_cotisdues.csv", "cotisation"},
		{"Sigfaible_pcoll.csv", "procol"},
		{"Sigfaible_etablissement_utf8.csv", "admin_urssaf"},
		{"Sigfaible_effectif_siret.csv", "effectif"},
		{"Sigfaible_effectif_siren.csv", "effectif_ent"},
		{"Sigfaible_delais.csv", "delai"},
		{"Sigfaible_ccsf.csv", "ccsf"},

		// guessed from dgefp files
		{"act_partielle_conso_depuis2014_FRANCE.csv", "apconso"},
		{"act_partielle_ddes_depuis2015_FRANCE.csv", "apdemande"},

		// others
		{"Sigfaibles_debits.csv", "debit"},
		{"Sigfaibles_debits2.csv", "debit"},
		{"diane_req_2002.csv", "diane"},
		{"diane_req_dom_2002.csv", "diane"},
		{"effectif_dom.csv", "effectif"},
		{"filter_siren_2002.csv", "filter"},
		{"sireneUL.csv", "sirene_ul"},
		{"StockEtablissement_utf8_geo.csv", "comptes"},
	}
	for _, testCase := range cases {
		t.Run("should return "+testCase.category+" for file "+testCase.name, func(t *testing.T) {
			got := GetFileType(testCase.name, DummyMetadataReader)
			assert.Equal(t, testCase.category, got)
		})
	}
}
