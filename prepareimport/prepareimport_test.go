package prepareimport

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io/ioutil"
	"os"
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
		expected := "could not find directory 1803/1803_01 in provided path"
		assert.Equal(t, expected, err.Error())
	})

	t.Run("Should warn if no filter is provided", func(t *testing.T) {
		dir := CreateTempFiles(t, dummyBatchKey, []string{"Sigfaibles_debits.csv"})
		_, err := PrepareImport(dir, dummyBatchKey, dummyDateFinEffectif)
		expected := "filter is missing: batch should include a filter or one effectif file"
		assert.Equal(t, expected, err.Error())
	})

	t.Run("Should warn if 2 effectif files are provided", func(t *testing.T) {
		dir := CreateTempFiles(t, dummyBatchKey, []string{"Sigfaible_effectif_siret.csv", "Sigfaible_effectif_siret2.csv"})
		_, err := PrepareImport(dir, dummyBatchKey, dummyDateFinEffectif)
		expected := "filter is missing: batch should include a filter or one effectif file"
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
		expected := FilesProperty{filter: {dummyBatchFile("filter_2002.csv")}}
		if assert.NoError(t, err) {
			assert.Equal(t, expected, res["files"])
		}
	})

	t.Run("Should detect and duplicate the filter file of the parent batch, given we are generating a sub-batch", func(t *testing.T) {
		// Set expectations
		subBatch := newSafeBatchKey("1803_01")
		parentBatch := newSafeBatchKey(subBatch.GetParentBatch())
		filterFile := newBatchFile(subBatch, "filter_siren_1803.csv")
		parentFilterFile := newBatchFile(parentBatch, "filter_siren_1803.csv")
		expectedFilesProp := FilesProperty{filter: {filterFile}}
		// Setup test environment
		parentDir := CreateTempFiles(t, parentBatch, []string{parentFilterFile.Name()})
		subBatchDir := filepath.Join(parentDir, parentBatch.String(), subBatch.String())
		os.Mkdir(subBatchDir, 0777)
		// Run the test
		res, err := PrepareImport(parentDir, subBatch, "2018-03-01")
		if assert.NoError(t, err) {
			assert.Equal(t, expectedFilesProp, res["files"])
			duplicatedFilePath := path.Join(parentDir, parentBatch.GetParentBatch(), filterFile.Path())
			assert.True(t, fileExists(duplicatedFilePath))
		}
	})

	t.Run("Should infer date_fin_effectif from the effectif file, given no filter was found", func(t *testing.T) {
		// Set expectations
		effectifFile := dummyBatchFile("Sigfaible_effectif_siret.csv")
		providedDateFinEffetif := ""
		expectedDateFinEffectif := NewDateFinEffectif(time.Date(2020, time.Month(1), 1, 0, 0, 0, 0, time.UTC)).MongoDate()
		// Setup test environment
		data, err := ioutil.ReadFile("../createfilter/test_data.csv")
		if err != nil {
			t.Fatal(err)
		}
		parentDir := CreateTempFilesWithContent(t, dummyBatchKey, map[string][]byte{
			effectifFile.Name(): data,
		})
		// Run the test
		res, err := PrepareImport(parentDir, dummyBatchKey, providedDateFinEffetif)
		if assert.NoError(t, err) {
			assert.Equal(t, expectedDateFinEffectif, res["param"].(ParamProperty).DateFinEffectif)
		}
	})

	t.Run("Should infer date_fin_effectif from the effectif file, given a filter was found", func(t *testing.T) {
		// Set expectations
		effectifFile := dummyBatchFile("Sigfaible_effectif_siret.csv")
		providedDateFinEffetif := ""
		expectedDateFinEffectif := NewDateFinEffectif(time.Date(2020, time.Month(1), 1, 0, 0, 0, 0, time.UTC)).MongoDate()
		// Setup test environment
		data, err := ioutil.ReadFile("../createfilter/test_data.csv")
		if err != nil {
			t.Fatal(err)
		}
		parentDir := CreateTempFilesWithContent(t, dummyBatchKey, map[string][]byte{
			effectifFile.Name(): data,
			"filter_2002.csv":   {},
		})
		// Run the test
		res, err := PrepareImport(parentDir, dummyBatchKey, providedDateFinEffetif)
		if assert.NoError(t, err) {
			assert.Equal(t, expectedDateFinEffectif, res["param"].(ParamProperty).DateFinEffectif)
		}
	})

	t.Run("Should infer the filter and date_fin_effectif from the effectif file of the parent batch, given we are generating a sub-batch", func(t *testing.T) {
		// Set expectations
		subBatch := newSafeBatchKey("1803_01")
		parentBatch := newSafeBatchKey(subBatch.GetParentBatch())
		filterFile := newBatchFile(subBatch, "filter_siren_1803.csv")
		parentEffectifFile := newBatchFile(parentBatch, "Sigfaible_effectif_siret.csv")
		expectedFilesProp := FilesProperty{filter: {filterFile}}
		expectedDateFinEffectif := NewDateFinEffectif(time.Date(2020, time.Month(1), 1, 0, 0, 0, 0, time.UTC)).MongoDate()
		// Setup test environment
		data, err := ioutil.ReadFile("../createfilter/test_data.csv")
		if err != nil {
			t.Fatal(err)
		}
		parentDir := CreateTempFilesWithContent(t, parentBatch, map[string][]byte{
			parentEffectifFile.Name(): data,
		})
		subBatchDir := filepath.Join(parentDir, parentBatch.String(), subBatch.String())
		os.Mkdir(subBatchDir, 0777)
		// Run the test
		res, err := PrepareImport(parentDir, subBatch, "")
		if assert.NoError(t, err) {
			assert.Equal(t, expectedFilesProp, res["files"])
			duplicatedFilePath := path.Join(parentDir, parentBatch.GetParentBatch(), filterFile.Path())
			assert.True(t, fileExists(duplicatedFilePath))
			assert.Equal(t, expectedDateFinEffectif, res["param"].(ParamProperty).DateFinEffectif)
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
			expected := []BatchFile{dummyBatchFile(testCase.id)}
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
		filterFileName := "filter_siren_" + dummyBatchKey.String() + ".csv"
		expected := FilesProperty{
			"effectif": {dummyBatchFile("Sigfaible_effectif_siret.csv")},
			"filter":   {dummyBatchFile(filterFileName)},
		}
		// check that the filter is listed in the "files" property
		if assert.NoError(t, err) {
			assert.Equal(t, expected, adminObject["files"])
		}
		// check that the filter file exists
		filterFilePath := path.Join(dir, dummyBatchKey.Path(), filterFileName)
		assert.True(t, fileExists(filterFilePath), "the filter file was not found: "+filterFilePath)
		// check that date_fin_effectif was detected from the effectif file
		validDateFinEffectif := time.Date(2020, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
		expectedDateFinEffectif := NewDateFinEffectif(validDateFinEffectif).MongoDate()
		actualDateFinEffectif := adminObject["param"].(ParamProperty).DateFinEffectif
		assert.Equal(t, expectedDateFinEffectif, actualDateFinEffectif)
	})

	t.Run("should create filter file even if effectif file is compressed", func(t *testing.T) {
		data, err := ioutil.ReadFile("../createfilter/test_data.csv")
		if err != nil {
			t.Fatal(err)
		}

		var compressedEffectifData bytes.Buffer
		zw := gzip.NewWriter(&compressedEffectifData)

		if _, err = zw.Write(data); err != nil {
			t.Fatal(err)
		}
		if err := zw.Close(); err != nil {
			t.Fatal(err)
		}

		// fileReader, err = gzip.NewReader(file)
		// if err != nil {
		// 	return file, nil, err
		// }

		metadata := `{ "MetaData": { "filename": "Sigfaible_effectif_siret.csv.gz", "goup-path": "acoss" }, "Size": 172391771 }`
		dir := CreateTempFilesWithContent(t, dummyBatchKey, map[string][]byte{
			// "Sigfaible_effectif_siret.csv.gz": compressedEffectifData.Bytes(),
			"719776012f6a124c3fab0f1c74fd585a":      {},
			"719776012f6a124c3fab0f1c74fd585a.info": []byte(metadata),
		})
		adminObject, err := PrepareImport(dir, dummyBatchKey, "") // => open /var/folders/v3/_c06yg_96tbf9kzmm0zq0y180000gn/T/example424586267/gzip:/1802/719776012f6a124c3fab0f1c74fd585a: no such file or directory
		filterFileName := "filter_siren_" + dummyBatchKey.String() + ".csv"
		expected := FilesProperty{
			"effectif": {dummyBatchFile("Sigfaible_effectif_siret.csv.gz")},
			"filter":   {dummyBatchFile(filterFileName)},
		}
		// check that the filter is listed in the "files" property
		if assert.NoError(t, err) {
			assert.Equal(t, expected, adminObject["files"])
		}
		// check that the filter file exists
		filterFilePath := path.Join(dir, dummyBatchKey.Path(), filterFileName)
		assert.True(t, fileExists(filterFilePath), "the filter file was not found: "+filterFilePath)
		// check that date_fin_effectif was detected from the effectif file
		validDateFinEffectif := time.Date(2020, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
		expectedDateFinEffectif := NewDateFinEffectif(validDateFinEffectif).MongoDate()
		actualDateFinEffectif := adminObject["param"].(ParamProperty).DateFinEffectif
		assert.Equal(t, expectedDateFinEffectif, actualDateFinEffectif)
	})
}
