package main

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper to create temporary files, and clean up after the execution of tests
func createTempFiles(t *testing.T, filenames []string) string {
	t.Helper()
	dir, err := ioutil.TempDir(os.TempDir(), "example")
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Cleanup(func() { os.RemoveAll(dir) })

	for _, filename := range filenames {
		tmpFilename := filepath.Join(dir, filename)
		if err := ioutil.WriteFile(tmpFilename, []byte{}, 0666); err != nil {
			t.Fatal(err.Error())
		}
	}

	return dir
}

func TestReadFilenames(t *testing.T) {
	t.Run("Should return filenames in a directory", func(t *testing.T) {
		dir := createTempFiles(t, []string{"tmpfile"})
		filenames, err := ReadFilenames(dir)
		if err != nil {
			t.Fatal(err.Error())
		}
		assert.Equal(t, []string{"tmpfile"}, filenames)
	})
}

func TestPrepareImport(t *testing.T) {
	t.Run("Should return a json with one file", func(t *testing.T) {
		dir := createTempFiles(t, []string{"Sigfaibles_debits.csv"})
		res, err := PrepareImport(dir)
		expected := AdminObject{
			"files": FilesProperty{"debit": []string{"Sigfaibles_debits.csv"}},
		}
		if assert.NoError(t, err) {
			assert.Equal(t, expected, res)
		}
	})

	cases := []struct {
		id       string
		filename string
		goupPath string
		filetype string
	}{
		{"9a047825d8173684b69994428449302f", "Sigfaible_debits.csv", "urssaf", "debit"},
		{"60d1bd320523904d8b8b427efbbd3928", "FICHIER_SF_2020_02.csv", "bdf", "bdf"},
	}

	for _, testCase := range cases {
		t.Run("Uploaded file originally named "+testCase.filename+" should be of type "+testCase.filetype, func(t *testing.T) {
			dir := createTempFiles(t, []string{testCase.id + ".bin"})

			tmpFilename := filepath.Join(dir, testCase.id+".info")
			content := []byte("{\"MetaData\":{\"filename\":\"" + testCase.filename + "\",\"goup-path\":\"" + testCase.goupPath + "\"}}")
			if err := ioutil.WriteFile(tmpFilename, content, 0666); err != nil {
				t.Fatal(err.Error())
			}

			res, err := PrepareImport(dir)
			expected := AdminObject{
				"files": FilesProperty{testCase.filetype: []string{testCase.id + ".bin"}},
			}
			if assert.NoError(t, err) {
				assert.Equal(t, expected, res)
			}
		})
	}

	t.Run("should return list of unsupported files", func(t *testing.T) {
		dir := createTempFiles(t, []string{"unsupported-file.csv"})
		_, err := PrepareImport(dir)
		var e *UnsupportedFilesError
		if assert.Error(t, err) && errors.As(err, &e) {
			assert.Equal(t, []string{"unsupported-file.csv"}, e.UnsupportedFiles)
		}
	})

	t.Run("should fail if missing .info file", func(t *testing.T) {
		dir := createTempFiles(t, []string{"lonely.bin"})
		assert.Panics(t, func() {
			PrepareImport(dir)
		})
	})
}

func TestPurePrepareImport(t *testing.T) {
	t.Run("Should return the filename in the debit property", func(t *testing.T) {
		filename := SimpleDataFile{"Sigfaibles_debits.csv"}

		res, err := PurePrepareImport([]DataFile{filename})
		expected := AdminObject{
			"files": FilesProperty{"debit": []string{"Sigfaibles_debits.csv"}},
		}
		if assert.NoError(t, err) {
			assert.Equal(t, expected, res)
		}
	})

	t.Run("Should return an empty json when there is no file", func(t *testing.T) {
		res, err := PurePrepareImport([]DataFile{})
		if assert.NoError(t, err) {
			assert.Equal(t, AdminObject{"files": FilesProperty{}}, res)
		}
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
		augmentedFiles := []DataFile{}
		for _, file := range files {
			augmentedFiles = append(augmentedFiles, SimpleDataFile{file})
		}
		res, err := PurePrepareImport(augmentedFiles)
		if assert.NoError(t, err) {
			resFilesProperty := res["files"].(FilesProperty)
			resultingFiles := []string{}
			for _, filenames := range resFilesProperty {
				resultingFiles = append(resultingFiles, filenames...)
			}
			assert.Subset(t, resultingFiles, files)
		}
	})
}

func TestPopulateFilesProperty(t *testing.T) {
	t.Run("PopulateFilesProperty should contain effectif file in \"effectif\" property", func(t *testing.T) {
		filesProperty, unsupportedFiles := PopulateFilesProperty([]DataFile{SimpleDataFile{"Sigfaibles_effectif_siret.csv"}})
		if assert.Len(t, unsupportedFiles, 0) {
			assert.Equal(t, []string{"Sigfaibles_effectif_siret.csv"}, filesProperty["effectif"])
		}
	})

	t.Run("PopulateFilesProperty should contain one debit file in \"debit\" property", func(t *testing.T) {
		filesProperty, unsupportedFiles := PopulateFilesProperty([]DataFile{SimpleDataFile{"Sigfaibles_debits.csv"}})
		if assert.Len(t, unsupportedFiles, 0) {
			assert.Equal(t, []string{"Sigfaibles_debits.csv"}, filesProperty["debit"])
		}
	})

	t.Run("PopulateFilesProperty should contain both debits files in \"debit\" property", func(t *testing.T) {
		filesProperty, unsupportedFiles := PopulateFilesProperty([]DataFile{SimpleDataFile{"Sigfaibles_debits.csv"}, SimpleDataFile{"Sigfaibles_debits2.csv"}})
		if assert.Len(t, unsupportedFiles, 0) {
			assert.Equal(t, []string{"Sigfaibles_debits.csv", "Sigfaibles_debits2.csv"}, filesProperty["debit"])
		}
	})

	t.Run("Should not include unsupported files", func(t *testing.T) {
		filesProperty, unsupportedFiles := PopulateFilesProperty([]DataFile{SimpleDataFile{"coco.csv"}})
		assert.Len(t, unsupportedFiles, 1)
		assert.Equal(t, FilesProperty{}, filesProperty)
	})
	t.Run("Should report unsupported files", func(t *testing.T) {
		_, unsupportedFiles := PopulateFilesProperty([]DataFile{SimpleDataFile{"coco.csv"}})
		assert.Equal(t, []string{"coco.csv"}, unsupportedFiles)
	})
}

func MakeMetadata(metadataFields MetadataProperty) UploadedFileMeta {
	return UploadedFileMeta{MetaData: metadataFields}
}

func TestExtractFileTypeFromMetadata(t *testing.T) {

	t.Run("should return \"debit\" for bin file which original name included \"debits\"", func(t *testing.T) {
		got := ExtractFileTypeFromMetadata("9a047825d8173684b69994428449302f.bin", MakeMetadata(MetadataProperty{
			"filename":  "Sigfaible_debits.csv",
			"goup-path": "urssaf",
		}))
		assert.Equal(t, "debit", got)
	})

	t.Run("should return \"bdf\" for bin file which came from bdf", func(t *testing.T) {
		got := ExtractFileTypeFromMetadata("60d1bd320523904d8b8b427efbbd3928.bin", MakeMetadata(MetadataProperty{
			"filename":  "FICHIER_SF_2020_02.csv",
			"goup-path": "bdf",
		}))
		assert.Equal(t, "bdf", got)
	})

	t.Run("should return \"interim\" for bin file which had a sas7dbat extension", func(t *testing.T) {
		got := ExtractFileTypeFromMetadata("ab8613ab66ebddb2db21e36b92fc5b70.bin", MakeMetadata(MetadataProperty{
			"filename":  "tab_19m10.sas7bdat",
			"goup-path": "dgefp",
		}))
		assert.Equal(t, "interim", got)
	})
}

func TestExtractFileTypeFromFilename(t *testing.T) {

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
		{"Diane_Export_4.txt", "diane"},
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
			got := ExtractFileTypeFromFilename(testCase.name)
			assert.Equal(t, testCase.category, got)
		})
	}
}
