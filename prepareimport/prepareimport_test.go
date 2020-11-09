package prepareimport

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

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
	t.Run("Should return a json with one file", func(t *testing.T) {
		dir := CreateTempFiles(t, dummyBatchKey, []string{"Sigfaibles_debits.csv"})
		res, err := PrepareImport(dir, dummyBatchKey, dummyDateFinEffectif)
		expected := FilesProperty{debit: []string{dummyBatchKey.Path() + "Sigfaibles_debits.csv"}}
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
		{"9a047825d8173684b69994428449302f", "Sigfaible_debits.csv", "urssaf", debit},
		{"60d1bd320523904d8b8b427efbbd3928", "FICHIER_SF_2020_02.csv", "bdf", bdf},
	}

	for _, testCase := range cases {
		t.Run("Uploaded file originally named "+testCase.filename+" should be of type "+string(testCase.filetype), func(t *testing.T) {

			dir := CreateTempFiles(t, dummyBatchKey, []string{testCase.id})

			tmpFilename := filepath.Join(dir, dummyBatchKey.String(), testCase.id+".info")
			content := []byte("{\"MetaData\":{\"filename\":\"" + dummyBatchKey.Path() + testCase.filename + "\",\"goup-path\":\"" + testCase.goupPath + "\"}}")
			if err := ioutil.WriteFile(tmpFilename, content, 0666); err != nil {
				t.Fatal(err.Error())
			}

			res, err := PrepareImport(dir, dummyBatchKey, dummyDateFinEffectif)
			expected := FilesProperty{testCase.filetype: []string{dummyBatchKey.Path() + testCase.id}}
			if assert.NoError(t, err) {
				assert.Equal(t, expected, res["files"])
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

	t.Run("should create filter file if a effectif files are present", func(t *testing.T) {
		data, err := ioutil.ReadFile("../createfilter/test_data.csv")
		if err != nil {
			t.Fatal(err)
		}
		dir := CreateTempFilesWithContent(t, dummyBatchKey, map[string][]byte{
			"Sigfaible_effectif_siret.csv": data,
		})
		res, err := PrepareImport(dir, dummyBatchKey, dummyDateFinEffectif)
		filterFileName := dummyBatchKey.Path() + "filter_siren_" + dummyBatchKey.String() + ".csv"
		expected := FilesProperty{
			"effectif": []string{dummyBatchKey.Path() + "Sigfaible_effectif_siret.csv"},
			"filter":   []string{filterFileName},
		}
		if assert.NoError(t, err) {
			assert.Equal(t, expected, res["files"])
		}
		filterFilePath := path.Join(dir, filterFileName)
		assert.True(t, fileExists(filterFilePath), "the filter file was not found: "+filterFilePath)
	})
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
