package prepareimport

import (
	"errors"
	"io/ioutil"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadFilenames(t *testing.T) {
	t.Run("Should return filenames in a directory", func(t *testing.T) {
		dir := CreateTempFiles(t, DUMMY_BATCHKEY, []string{"tmpfile"})
		filenames, err := ReadFilenames(path.Join(dir, DUMMY_BATCHKEY.String()))
		if err != nil {
			t.Fatal(err.Error())
		}
		assert.Equal(t, []string{"tmpfile"}, filenames)
	})
}

func TestPrepareImport(t *testing.T) {
	t.Run("Should return a json with one file", func(t *testing.T) {
		dir := CreateTempFiles(t, DUMMY_BATCHKEY, []string{"Sigfaibles_debits.csv"})
		res, err := PrepareImport(dir, DUMMY_BATCHKEY, DUMMY_DATE_FIN_EFFECTIF)
		expected := FilesProperty{DEBIT: []string{DUMMY_BATCHKEY.Path() + "Sigfaibles_debits.csv"}}
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

			dir := CreateTempFiles(t, DUMMY_BATCHKEY, []string{testCase.id})

			tmpFilename := filepath.Join(dir, DUMMY_BATCHKEY.String(), testCase.id+".info")
			content := []byte("{\"MetaData\":{\"filename\":\"" + DUMMY_BATCHKEY.Path() + testCase.filename + "\",\"goup-path\":\"" + testCase.goupPath + "\"}}")
			if err := ioutil.WriteFile(tmpFilename, content, 0666); err != nil {
				t.Fatal(err.Error())
			}

			res, err := PrepareImport(dir, DUMMY_BATCHKEY, DUMMY_DATE_FIN_EFFECTIF)
			expected := FilesProperty{testCase.filetype: []string{DUMMY_BATCHKEY.Path() + testCase.id}}
			if assert.NoError(t, err) {
				assert.Equal(t, expected, res["files"])
			}
		})
	}

	t.Run("should return list of unsupported files", func(t *testing.T) {
		dir := CreateTempFiles(t, DUMMY_BATCHKEY, []string{"unsupported-file.csv"})
		_, err := PrepareImport(dir, DUMMY_BATCHKEY, DUMMY_DATE_FIN_EFFECTIF)
		var e *UnsupportedFilesError
		if assert.Error(t, err) && errors.As(err, &e) {
			assert.Equal(t, []string{DUMMY_BATCHKEY.Path() + "unsupported-file.csv"}, e.UnsupportedFiles)
		}
	})

	t.Run("should fail if missing .info file", func(t *testing.T) {
		dir := CreateTempFiles(t, DUMMY_BATCHKEY, []string{"lonely"})
		assert.Panics(t, func() {
			PrepareImport(dir, DUMMY_BATCHKEY, DUMMY_DATE_FIN_EFFECTIF)
		})
	})
}
