package main

import (
	"bytes"
	"log"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	t.Run("prepare-import golden file", func(t *testing.T) {

		dir := createTempFiles(t, "Sigfaibles_effectif_siret.csv")

		cmd := exec.Command("./prepare-import", "--path", dir)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()

		if err != nil {
			log.Fatal(err)
		}
		expected := "{\n  \"files\": {\n    \"effectif\": [\n      \"Sigfaibles_effectif_siret.csv\"\n    ]\n  }\n}\n"
		assert.Equal(t, expected, out.String())
	})
}
