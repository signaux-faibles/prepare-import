package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		expected := "{\n  \"files\": {\n    \"debit\": [\n      \"Sigfaibles_debits.csv\"\n    ],\n    \"effectif\": [\n      \"Sigfaibles_effectif_siret.csv\"\n    ]\n  }\n}\n"
		assert.Equal(t, expected, out.String())
	})
}
