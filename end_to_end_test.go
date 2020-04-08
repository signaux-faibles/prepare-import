package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var outGoldenFile = "end_to_end_golden.txt"
var errGoldenFile = "end_to_end_golden_err.txt"

var updateGoldenFile = flag.Bool("update", false, "Update the expected test values in golden file")

func diffWithGoldenFile(filename string, updateGoldenFile bool, cmdOutput bytes.Buffer) []byte {

	if updateGoldenFile {
		ioutil.WriteFile(filename, cmdOutput.Bytes(), 0644)
	}
	expected, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return expected
}

func TestMain(t *testing.T) {
	t.Run("prepare-import golden file", func(t *testing.T) {

		dir := createTempFiles(t, []string{"Sigfaibles_effectif_siret.csv", "Sigfaibles_debits.csv", "abcdef.bin", "unsupported.csv"})

		content := []byte("{\"MetaData\":{\"filename\":\"FICHIER_SF_2020_02.csv\",\"goup-path\":\"bdf\"}}")
		ioutil.WriteFile(filepath.Join(dir, "abcdef.info"), content, 0644)

		cmd := exec.Command("./prepare-import", "--path", dir)
		var cmdOutput bytes.Buffer
		var cmdError bytes.Buffer
		cmd.Stdout = &cmdOutput
		cmd.Stderr = &cmdError
		err := cmd.Run()
		fmt.Printf("stderr: %q\n", cmdError.String())
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("stdout: %q\n", cmdOutput.String())

		expectedOutput := diffWithGoldenFile(outGoldenFile, *updateGoldenFile, cmdOutput)
		expectedError := diffWithGoldenFile(errGoldenFile, *updateGoldenFile, cmdError)

		assert.Equal(t, string(expectedOutput), cmdOutput.String())
		assert.Equal(t, string(expectedError), cmdError.String())
	})
}
