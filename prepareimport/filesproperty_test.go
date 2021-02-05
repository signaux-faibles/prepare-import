package prepareimport

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPopulateFilesProperty(t *testing.T) {
	t.Run("Should return an empty json when there is no file", func(t *testing.T) {
		filesProperty, unsupportedFiles := PopulateFilesPropertyFromDataFiles([]DataFile{}, dummyBatchKey)
		assert.Len(t, unsupportedFiles, 0)
		assert.Equal(t, FilesProperty{}, filesProperty)
	})

	t.Run("PopulateFilesProperty should contain effectif file in \"effectif\" property", func(t *testing.T) {
		filesProperty, unsupportedFiles := PopulateFilesPropertyFromDataFiles([]DataFile{SimpleDataFile{"Sigfaibles_effectif_siret.csv"}}, dummyBatchKey)
		if assert.Len(t, unsupportedFiles, 0) {
			assert.Equal(t, []BatchFile{dummyBatchFile("Sigfaibles_effectif_siret.csv")}, filesProperty[effectif])
		}
	})

	t.Run("PopulateFilesProperty should contain one debit file in \"debit\" property", func(t *testing.T) {
		filesProperty, unsupportedFiles := PopulateFilesPropertyFromDataFiles([]DataFile{SimpleDataFile{"Sigfaibles_debits.csv"}}, dummyBatchKey)
		expected := FilesProperty{debit: {dummyBatchFile("Sigfaibles_debits.csv")}}
		assert.Len(t, unsupportedFiles, 0)
		assert.Equal(t, expected, filesProperty)
	})

	t.Run("PopulateFilesProperty should contain both debits files in \"debit\" property", func(t *testing.T) {
		filesProperty, unsupportedFiles := PopulateFilesPropertyFromDataFiles([]DataFile{SimpleDataFile{"Sigfaibles_debits.csv"}, SimpleDataFile{"Sigfaibles_debits2.csv"}}, dummyBatchKey)
		if assert.Len(t, unsupportedFiles, 0) {
			assert.Equal(t, []BatchFile{dummyBatchFile("Sigfaibles_debits.csv"), dummyBatchFile("Sigfaibles_debits2.csv")}, filesProperty[debit])
		}
	})

	t.Run("Should support multiple types of csv files", func(t *testing.T) {
		type File struct {
			Type     ValidFileType
			Filename string
		}
		files := []File{
			{"diane", "diane_req_2002.csv"},               // --> DIANE
			{"diane", "diane_req_dom_2002.csv"},           // --> DIANE
			{"effectif", "effectif_dom.csv"},              // --> EFFECTIF
			{"filter", "filter_siren_2002.csv"},           // --> FILTER
			{"sirene_ul", "sireneUL.csv"},                 // --> SIRENE_UL
			{"sirene", "StockEtablissement_utf8_geo.csv"}, // --> SIRENE
		}
		expectedFiles := FilesProperty{}
		inputFiles := []DataFile{}
		for _, file := range files {
			expectedFiles[file.Type] = append(expectedFiles[file.Type], dummyBatchFile(file.Filename))
			inputFiles = append(inputFiles, SimpleDataFile{file.Filename})
		}
		resFilesProperty, unsupportedFiles := PopulateFilesPropertyFromDataFiles(inputFiles, dummyBatchKey)
		assert.Len(t, unsupportedFiles, 0)
		assert.Equal(t, expectedFiles, resFilesProperty)
	})

	t.Run("Should not include unsupported files", func(t *testing.T) {
		filesProperty, unsupportedFiles := PopulateFilesPropertyFromDataFiles([]DataFile{SimpleDataFile{"coco.csv"}}, dummyBatchKey)
		assert.Len(t, unsupportedFiles, 1)
		assert.Equal(t, FilesProperty{}, filesProperty)
	})

	t.Run("Should report unsupported files", func(t *testing.T) {
		_, unsupportedFiles := PopulateFilesPropertyFromDataFiles([]DataFile{SimpleDataFile{"coco.csv"}}, dummyBatchKey)
		assert.Equal(t, []string{dummyBatchKey.Path() + "coco.csv"}, unsupportedFiles)
	})

	t.Run("Should skip subdirectories", func(t *testing.T) {
		subBatch := newSafeBatchKey("1803_01")
		parentBatch := subBatch.GetParentBatch()
		parentDir := CreateTempFiles(t, newSafeBatchKey(parentBatch), []string{})
		subBatchDir := filepath.Join(parentDir, parentBatch, subBatch.String())
		os.Mkdir(subBatchDir, 0777)
		parentFilesProperty, unsupportedFiles := PopulateFilesProperty(parentDir, newSafeBatchKey(parentBatch))
		assert.Equal(t, []string{}, unsupportedFiles)
		assert.Equal(t, FilesProperty{}, parentFilesProperty)
	})

	t.Run("Should forward the size of a gzipped file provided with metadata", func(t *testing.T) {

		metadata := `{ "MetaData": { "filename": "Sigfaibles_debits.csv", "goup-path": "" }, "Size": 254781489 }` // thresholdPerGzippedFileType["debit"]
		dir := CreateTempFilesWithContent(t, dummyBatchKey, map[string][]byte{
			"083fe617e80f2e30a21598d38a854bc6":      {},
			"083fe617e80f2e30a21598d38a854bc6.info": []byte(metadata),
		})

		// res, err := PrepareImport(dir, dummyBatchKey, dummyDateFinEffectif)
		// expected := []BatchFile{dummyBatchFile(testCase.id)}
		// if assert.NoError(t, err) {
		// 	assert.Equal(t, expected, res["files"].(FilesProperty)[testCase.filetype])
		// }

		// inputFiles := []DataFile{
		// 	UploadedDataFile{
		// 		filename: "Sigfaibles_debits.csv",
		// 		path:     ".",
		// 	},
		// }
		expectedFiles := FilesProperty{"debit": {
			batchFile{
				batchKey:    dummyBatchKey,
				filename:    "Sigfaibles_debits.csv",
				gzippedSize: 254781489, // thresholdPerGzippedFileType["debit"]
			},
		}}
		resFilesProperty, unsupportedFiles := PopulateFilesProperty(dir, dummyBatchKey)
		assert.Len(t, unsupportedFiles, 0)
		assert.Equal(t, expectedFiles, resFilesProperty)
	})
}
