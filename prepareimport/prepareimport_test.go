package prepareimport

import (
	"errors"
	"io/ioutil"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReadFilenames(t *testing.T) {
	t.Run("Should return filenames in a directory", func(t *testing.T) {
		dir := CreateTempFiles(t, dummyBatchKey, []string{"tmpfile"})
		filenames, err := ReadFilenames(path.Join(dir, dummyBatchKey.String()))
		if err != nil {
			t.Fatal(err.Error())
		}
		assert.Equal(t, []string{"tmpfile"}, filenames)
	})
}

func TestPrepareImport(t *testing.T) {
	t.Run("Should warn if the batch was not found in the specified directory", func(t *testing.T) {
		wantedBatch := newSafeBatchKey("1803") // different of dummyBatchKey
		parentDir := CreateTempFiles(t, dummyBatchKey, []string{})
		_, err := PrepareImport(parentDir, wantedBatch, "")
		expected := "could not find directory 1803 in provided path"
		assert.Equal(t, expected, err.Error())
	})

	t.Run("Should warn if the sub-batch was not found in the specified directory", func(t *testing.T) {
		subBatch := newSafeBatchKey("1803_01")
		parentBatch := newSafeBatchKey("1803")
		parentDir := CreateTempFiles(t, parentBatch, []string{})
		_, err := PrepareImport(parentDir, subBatch, "")
		expected := "could not find directory 1803/01 in provided path"
		assert.Equal(t, expected, err.Error())
	})

	t.Run("Should warn if no filter is provided", func(t *testing.T) {
		dir := CreateTempFiles(t, dummyBatchKey, []string{"Sigfaibles_debits.csv"})
		_, err := PrepareImport(dir, dummyBatchKey, dummyDateFinEffectif)
		expected := "filter is missing: please include a filter or an effectif file"
		assert.Equal(t, expected, err.Error())
	})

	t.Run("Should warn if 2 effectif files are provided", func(t *testing.T) {
		dir := CreateTempFiles(t, dummyBatchKey, []string{"Sigfaible_effectif_siret.csv", "Sigfaible_effectif_siret2.csv"})
		_, err := PrepareImport(dir, dummyBatchKey, dummyDateFinEffectif)
		expected := "generating a filter requires just 1 effectif file, found: [/1802/Sigfaible_effectif_siret.csv /1802/Sigfaible_effectif_siret2.csv]"
		assert.Equal(t, expected, err.Error())
	})

	t.Run("Should warn if neither effectif and date_fin_effectif are provided", func(t *testing.T) {
		dir := CreateTempFiles(t, dummyBatchKey, []string{"filter_2002.csv"})
		_, err := PrepareImport(dir, dummyBatchKey, "")
		expected := "date_fin_effectif is missing or invalid: "
		assert.Equal(t, expected, err.Error())
	})

	t.Run("Should return a json with one file", func(t *testing.T) {
		dir := CreateTempFiles(t, dummyBatchKey, []string{"filter_2002.csv"})
		res, err := PrepareImport(dir, dummyBatchKey, dummyDateFinEffectif)
		expected := FilesProperty{filter: []string{dummyBatchKey.Path() + "filter_2002.csv"}}
		if assert.NoError(t, err) {
			assert.Equal(t, expected, res["files"])
		}
	})

	t.Run("Should return an _id property", func(t *testing.T) {
		batch := newSafeBatchKey("1802")
		dir := CreateTempFiles(t, batch, []string{"filter_2002.csv"})
		res, err := PrepareImport(dir, batch, dummyDateFinEffectif)
		if assert.NoError(t, err) {
			assert.Equal(t, IDProperty{batch, "batch"}, res["_id"])
		}
	})

	cases := []struct {
		id       string
		filename string
		goupPath string
		filetype ValidFileType
	}{
		{"9a047825d8173684b69994428449302f", "Sigfaible_debits.csv", "urssaf", debit},
		{"60d1bd320523904d8b8b427efbbd3928", "FICHIER_SF_2020_02.csv", "bdf", bdf},
	}

	for _, testCase := range cases {
		t.Run("Uploaded file originally named "+testCase.filename+" should be of type "+string(testCase.filetype), func(t *testing.T) {

			dir := CreateTempFiles(t, dummyBatchKey, []string{testCase.id, "filter_2002.csv"})

			tmpFilename := filepath.Join(dir, dummyBatchKey.String(), testCase.id+".info")
			content := []byte("{\"MetaData\":{\"filename\":\"" + dummyBatchKey.Path() + testCase.filename + "\",\"goup-path\":\"" + testCase.goupPath + "\"}}")
			if err := ioutil.WriteFile(tmpFilename, content, 0666); err != nil {
				t.Fatal(err.Error())
			}

			res, err := PrepareImport(dir, dummyBatchKey, dummyDateFinEffectif)
			expected := []string{dummyBatchKey.Path() + testCase.id}
			if assert.NoError(t, err) {
				assert.Equal(t, expected, res["files"].(FilesProperty)[testCase.filetype])
			}
		})
	}

	t.Run("should return list of unsupported files", func(t *testing.T) {
		dir := CreateTempFiles(t, dummyBatchKey, []string{"unsupported-file.csv"})
		_, err := PrepareImport(dir, dummyBatchKey, dummyDateFinEffectif)
		var e *UnsupportedFilesError
		if assert.Error(t, err) && errors.As(err, &e) {
			assert.Equal(t, []string{dummyBatchKey.Path() + "unsupported-file.csv"}, e.UnsupportedFiles)
		}
	})

	t.Run("should fail if missing .info file", func(t *testing.T) {
		dir := CreateTempFiles(t, dummyBatchKey, []string{"lonely"})
		assert.Panics(t, func() {
			PrepareImport(dir, dummyBatchKey, dummyDateFinEffectif)
		})
	})

	t.Run("should create filter file and fill date_fin_effectif if an effectif file is present", func(t *testing.T) {
		data, err := ioutil.ReadFile("../createfilter/test_data.csv")
		if err != nil {
			t.Fatal(err)
		}
		dir := CreateTempFilesWithContent(t, dummyBatchKey, map[string][]byte{
			"Sigfaible_effectif_siret.csv": data,
		})
		adminObject, err := PrepareImport(dir, dummyBatchKey, "")
		filterFileName := dummyBatchKey.Path() + "filter_siren_" + dummyBatchKey.String() + ".csv"
		expected := FilesProperty{
			"effectif": []string{dummyBatchKey.Path() + "Sigfaible_effectif_siret.csv"},
			"filter":   []string{filterFileName},
		}
		// check that the filter is listed in the "files" property
		if assert.NoError(t, err) {
			assert.Equal(t, expected, adminObject["files"])
		}
		// check that the filter file exists
		filterFilePath := path.Join(dir, filterFileName)
		assert.True(t, fileExists(filterFilePath), "the filter file was not found: "+filterFilePath)
		// check that date_fin_effectif was detected from the effectif file
		validDateFinEffectif := time.Date(2020, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
		expectedDateFinEffectif := NewDateFinEffectif(validDateFinEffectif).MongoDate()
		actualDateFinEffectif := adminObject["param"].(ParamProperty).DateFinEffectif
		assert.Equal(t, expectedDateFinEffectif, actualDateFinEffectif)
	})
}
