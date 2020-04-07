package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"testing"
)

func TestMain(t *testing.T) {
	t.Run("prepare-import golden file", func(t *testing.T) {

		// dir := createTempFiles(t, "tmpfile.csv")
		dir := createTempFiles(t, "Sigfaibles_effectif_siret.csv")

		cmd := exec.Command("prepare-import", "--path", dir)
		cmd.Stdin = strings.NewReader("some input")
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()

		fmt.Printf("in all caps: %q\n", out.String())
		if err != nil {
			log.Fatal(err)
		}
		// https://golang.org/pkg/os/exec/#Command

		// filesProperty, unsupportedFiles := PopulateFilesProperty([]DataFile{SimpleDataFile{"Sigfaibles_effectif_siret.csv"}})
		// if assert.Len(t, unsupportedFiles, 0) {
		// 	assert.Equal(t, []string{"Sigfaibles_effectif_siret.csv"}, filesProperty["effectif"])
		// }
	})
}
