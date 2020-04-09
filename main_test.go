package main

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const DUMMY_BATCHKEY batchKeyType = "1802"

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
		res, err := PrepareImport(dir, DUMMY_BATCHKEY)
		expected := FilesProperty{DEBIT: []string{"Sigfaibles_debits.csv"}}
		if assert.NoError(t, err) {
			assert.Equal(t, expected, res["files"])
		}
	})

	cases := []struct {
		id       string
		filename string
		goupPath string
		filetype ValidFileType
	}{
		{"9a047825d8173684b69994428449302f", "Sigfaible_debits.csv", "urssaf", DEBIT},
		{"60d1bd320523904d8b8b427efbbd3928", "FICHIER_SF_2020_02.csv", "bdf", BDF},
	}

	for _, testCase := range cases {
		t.Run("Uploaded file originally named "+testCase.filename+" should be of type "+string(testCase.filetype), func(t *testing.T) {
			dir := createTempFiles(t, []string{testCase.id + ".bin"})

			tmpFilename := filepath.Join(dir, testCase.id+".info")
			content := []byte("{\"MetaData\":{\"filename\":\"" + testCase.filename + "\",\"goup-path\":\"" + testCase.goupPath + "\"}}")
			if err := ioutil.WriteFile(tmpFilename, content, 0666); err != nil {
				t.Fatal(err.Error())
			}

			res, err := PrepareImport(dir, DUMMY_BATCHKEY)
			expected := FilesProperty{testCase.filetype: []string{testCase.id + ".bin"}}
			if assert.NoError(t, err) {
				assert.Equal(t, expected, res["files"])
			}
		})
	}

	t.Run("should return list of unsupported files", func(t *testing.T) {
		dir := createTempFiles(t, []string{"unsupported-file.csv"})
		_, err := PrepareImport(dir, DUMMY_BATCHKEY)
		var e *UnsupportedFilesError
		if assert.Error(t, err) && errors.As(err, &e) {
			assert.Equal(t, []string{"unsupported-file.csv"}, e.UnsupportedFiles)
		}
	})

	t.Run("should fail if missing .info file", func(t *testing.T) {
		dir := createTempFiles(t, []string{"lonely.bin"})
		assert.Panics(t, func() {
			PrepareImport(dir, DUMMY_BATCHKEY)
		})
	})
}

func TestBatchKey(t *testing.T) {

	t.Run("Should accept valid batch key", func(t *testing.T) {
		_, err := BatchKey("1802")
		assert.NoError(t, err)
	})

	t.Run("Should fail if batch key is invalid", func(t *testing.T) {
		_, err := BatchKey("")
		assert.Error(t, err, "la clé du batch doit respecter le format requis AAMM")
	})
}

func TestPopulateAdminObject(t *testing.T) {
	t.Run("Should return the filename in the debit property", func(t *testing.T) {
		filename := SimpleDataFile{"Sigfaibles_debits.csv"}

		res, err := PopulateAdminObject([]DataFile{filename}, DUMMY_BATCHKEY)
		expected := FilesProperty{DEBIT: []string{"Sigfaibles_debits.csv"}}
		if assert.NoError(t, err) {
			assert.Equal(t, expected, res["files"])
		}
	})

	t.Run("Should return an empty complete_types property", func(t *testing.T) {
		filename := SimpleDataFile{"Sigfaibles_debits.csv"}

		res, err := PopulateAdminObject([]DataFile{filename}, DUMMY_BATCHKEY)
		expected := []ValidFileType{}
		if assert.NoError(t, err) {
			assert.Equal(t, expected, res["complete_types"])
		}
	})

	t.Run("Should return apconso as a complete_type", func(t *testing.T) {
		filename := SimpleDataFile{"act_partielle_conso_depuis2014_FRANCE.csv"}
		res, err := PopulateAdminObject([]DataFile{filename}, DUMMY_BATCHKEY)
		expected := []ValidFileType{APCONSO}
		if assert.NoError(t, err) {
			assert.Equal(t, expected, res["complete_types"])
		}
	})

	t.Run("Should return an empty json when there is no file", func(t *testing.T) {
		res, err := PopulateAdminObject([]DataFile{}, DUMMY_BATCHKEY)
		if assert.NoError(t, err) {
			assert.Equal(t, FilesProperty{}, res["files"])
		}
	})

	t.Run("Should support multiple types of csv files", func(t *testing.T) {
		files := []string{
			"diane_req_2002.csv",              // --> DIANE
			"diane_req_dom_2002.csv",          // --> DIANE
			"effectif_dom.csv",                // --> EFFECTIF
			"filter_siren_2002.csv",           // --> FILTER
			"sireneUL.csv",                    // --> SIRENE_UL
			"StockEtablissement_utf8_geo.csv", // --> SIRENE
		}
		augmentedFiles := []DataFile{}
		for _, file := range files {
			augmentedFiles = append(augmentedFiles, SimpleDataFile{file})
		}
		res, err := PopulateAdminObject(augmentedFiles, DUMMY_BATCHKEY)
		if assert.NoError(t, err) {
			resFilesProperty := res["files"].(FilesProperty)
			resultingFiles := []string{}
			for _, filenames := range resFilesProperty {
				resultingFiles = append(resultingFiles, filenames...)
			}
			assert.Subset(t, resultingFiles, files)
		}
	})

	t.Run("Should return an _id property", func(t *testing.T) {
		res, err := PopulateAdminObject([]DataFile{}, "1802")
		if assert.NoError(t, err) {
			assert.Equal(t, IDProperty{"1802", "batch"}, res["_id"])
		}
	})

	t.Run("Should return a date_fin consistent with batch key", func(t *testing.T) {
		res, err := PopulateAdminObject([]DataFile{}, "1912") // ~= 12/2019
		expected := map[string]string{"$date": "2019-12-01T00:00:00.000+0000"}
		if assert.NoError(t, err) {
			assert.Equal(t, expected, res["param"].(map[string]map[string]string)["date_fin"])
		}
	})
}

func TestPopulateFilesProperty(t *testing.T) {
	t.Run("PopulateFilesProperty should contain effectif file in \"effectif\" property", func(t *testing.T) {
		filesProperty, unsupportedFiles := PopulateFilesProperty([]DataFile{SimpleDataFile{"Sigfaibles_effectif_siret.csv"}})
		if assert.Len(t, unsupportedFiles, 0) {
			assert.Equal(t, []string{"Sigfaibles_effectif_siret.csv"}, filesProperty[EFFECTIF])
		}
	})

	t.Run("PopulateFilesProperty should contain one debit file in \"debit\" property", func(t *testing.T) {
		filesProperty, unsupportedFiles := PopulateFilesProperty([]DataFile{SimpleDataFile{"Sigfaibles_debits.csv"}})
		if assert.Len(t, unsupportedFiles, 0) {
			assert.Equal(t, []string{"Sigfaibles_debits.csv"}, filesProperty[DEBIT])
		}
	})

	t.Run("PopulateFilesProperty should contain both debits files in \"debit\" property", func(t *testing.T) {
		filesProperty, unsupportedFiles := PopulateFilesProperty([]DataFile{SimpleDataFile{"Sigfaibles_debits.csv"}, SimpleDataFile{"Sigfaibles_debits2.csv"}})
		if assert.Len(t, unsupportedFiles, 0) {
			assert.Equal(t, []string{"Sigfaibles_debits.csv", "Sigfaibles_debits2.csv"}, filesProperty[DEBIT])
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
		assert.Equal(t, DEBIT, got)
	})

	t.Run("should return \"bdf\" for bin file which came from bdf", func(t *testing.T) {
		got := ExtractFileTypeFromMetadata("60d1bd320523904d8b8b427efbbd3928.bin", MakeMetadata(MetadataProperty{
			"filename":  "FICHIER_SF_2020_02.csv",
			"goup-path": "bdf",
		}))
		assert.Equal(t, BDF, got)
	})

	t.Run("should return \"interim\" for bin file which had a sas7dbat extension", func(t *testing.T) {
		got := ExtractFileTypeFromMetadata("ab8613ab66ebddb2db21e36b92fc5b70.bin", MakeMetadata(MetadataProperty{
			"filename":  "tab_19m10.sas7bdat",
			"goup-path": "dgefp",
		}))
		assert.Equal(t, INTERIM, got)
	})
}

func TestExtractFileTypeFromFilename(t *testing.T) {

	// inspired by https://github.com/golang/go/wiki/TableDrivenTests
	cases := []struct {
		name     string
		category ValidFileType
	}{
		// guessed from urssaf files found on stockage/goub server
		{"Sigfaible_debits.csv", DEBIT},
		{"Sigfaible_cotisdues.csv", COTISATION},
		{"Sigfaible_pcoll.csv", PROCOL},
		{"Sigfaible_etablissement_utf8.csv", ADMIN_URSSAF},
		{"Sigfaible_effectif_siret.csv", EFFECTIF},
		{"Sigfaible_effectif_siren.csv", EFFECTIF_ENT},
		{"Sigfaible_delais.csv", DELAI},
		{"Sigfaible_ccsf.csv", CCSF},

		// guessed from dgefp files
		{"act_partielle_conso_depuis2014_FRANCE.csv", APCONSO},
		{"act_partielle_ddes_depuis2015_FRANCE.csv", APDEMANDE},

		// others
		{"Diane_Export_4.txt", DIANE},
		{"Sigfaibles_debits.csv", DEBIT},
		{"Sigfaibles_debits2.csv", DEBIT},
		{"diane_req_2002.csv", DIANE},
		{"diane_req_dom_2002.csv", DIANE},
		{"effectif_dom.csv", EFFECTIF},
		{"filter_siren_2002.csv", FILTER},
		{"sireneUL.csv", SIRENE_UL},
		{"StockEtablissement_utf8_geo.csv", SIRENE},
	}
	for _, testCase := range cases {
		t.Run("should return "+string(testCase.category)+" for file "+testCase.name, func(t *testing.T) {
			got := ExtractFileTypeFromFilename(testCase.name)
			assert.Equal(t, testCase.category, got)
		})
	}
}
