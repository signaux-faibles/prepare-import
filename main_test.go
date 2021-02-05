package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"os/exec"
	"testing"

	"github.com/signaux-faibles/prepare-import/createfilter"
	"github.com/signaux-faibles/prepare-import/prepareimport"
	"github.com/stretchr/testify/assert"
)

var outGoldenFile = "end_to_end_golden.txt"
var errGoldenFile = "end_to_end_golden_err.txt"

var updateGoldenFile = flag.Bool("update", false, "Update the expected test values in golden file")

func TestMain(t *testing.T) {
	t.Run("prepare-import golden file", func(t *testing.T) {

		batch := "1802"

		batchKey, _ := prepareimport.NewBatchKey(batch)

		effectifData, err := ioutil.ReadFile("./createfilter/test_data.csv")
		if err != nil {
			t.Fatal(err)
		}
		parentDir := prepareimport.CreateTempFilesWithContent(t, batchKey, map[string][]byte{
			"Sigfaibles_effectif_siret.csv": effectifData,
			"Sigfaibles_debits.csv":         {},
			"abcdef":                        {},
			"abcdef.info":                   []byte(`{ "MetaData": { "filename": "FICHIER_SF_2020_02.csv", "goup-path": "bdf" } }`),
			"unsupported.csv":               {},
			"E_202011095813_Retro-Paydex_20201207.csv": {},
			// "083fe617e80f2e30a21598d38a854bc6":         {},
			// "083fe617e80f2e30a21598d38a854bc6.info":    []byte(`{ "MetaData": { "filename": "Sigfaible_pcoll.csv", "goup-path": "" } }`),
		})

		cmds := []*exec.Cmd{
			exec.Command(
				"./prepare-import",
				"--path", parentDir,
				"--batch", batch,
				"--date-fin-effectif", "2014-01-01",
			), // paramètres valides
			exec.Command("./prepare-import", "--path", parentDir, "--batch", "180"), // nom de batch invalide
		}
		var cmdOutput bytes.Buffer
		var cmdError bytes.Buffer
		for _, cmd := range cmds {
			cmd.Stdout = &cmdOutput
			cmd.Stderr = &cmdError
			err := cmd.Run()
			if err != nil {
				// log.Fatal(err)
			}
		}

		expectedOutput := createfilter.DiffWithGoldenFile(outGoldenFile, *updateGoldenFile, cmdOutput)
		expectedError := createfilter.DiffWithGoldenFile(errGoldenFile, *updateGoldenFile, cmdError)

		assert.Equal(t, string(expectedOutput), cmdOutput.String())
		assert.Equal(t, string(expectedError), cmdError.String())
	})
}
