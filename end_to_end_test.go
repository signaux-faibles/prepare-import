package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

var goldenFile = "end_to_end_golden.txt"

var update = flag.Bool("update", false, "Update the expected test values in golden file")

func TestMain(t *testing.T) {
	t.Run("prepare-import golden file", func(t *testing.T) {

		dir := createTempFiles(t, []string{"Sigfaibles_effectif_siret.csv", "Sigfaibles_debits.csv"})

		cmd := exec.Command("./prepare-import", "--path", dir)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()

		fmt.Printf("stdout: %q\n", out.String())
		if err != nil {
			log.Fatal(err)
		}
		if *update {
			ioutil.WriteFile(goldenFile, out.Bytes(), 0644)
		}
		expected, err := ioutil.ReadFile(goldenFile)
		if err != nil {
			log.Fatal(err)
		}
		// expected := "{\n  \"files\": {\n    \"debit\": [\n      \"Sigfaibles_debits.csv\"\n    ],\n    \"effectif\": [\n      \"Sigfaibles_effectif_siret.csv\"\n    ]\n  }\n}\n"
		assert.Equal(t, expected, out.Bytes())
	})
}
