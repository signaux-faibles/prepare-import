package tools

import (
	"bytes"
	"flag"
	"os/exec"
	"testing"

	"github.com/signaux-faibles/prepare-import/createfilter"
	"github.com/stretchr/testify/assert"
)

var outGoldenFile = "filter_to_diane_golden.txt"

var updateGoldenFile = flag.Bool("update", false, "Update the expected test values in golden file")

func TestFilterToDiane(t *testing.T) {
	t.Run("filter_to_diane golden file", func(t *testing.T) {

		cmd := exec.Command("./filter_to_diane.awk", "-v", "var_num=CF00012", "filter_to_diane_testdata.txt")
		var cmdOutput bytes.Buffer
		cmd.Stdout = &cmdOutput
		err := cmd.Run()
		if err != nil {
			t.Fatal(err)
		}

		expectedOutput := createfilter.DiffWithGoldenFile(outGoldenFile, *updateGoldenFile, cmdOutput)

		assert.Equal(t, string(expectedOutput), cmdOutput.String())
	})
}
