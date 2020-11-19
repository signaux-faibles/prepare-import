package tools

import (
	"bytes"
	"flag"
	"os"
	"os/exec"
	"testing"

	"github.com/signaux-faibles/prepare-import/createfilter"
	"github.com/stretchr/testify/assert"
)

var updateGoldenFile = flag.Bool("update", false, "Update the expected test values in golden file")

func TestFilterToDiane(t *testing.T) {
	t.Run("filter_to_diane golden file", func(t *testing.T) {

		var outGoldenFile = "filter_to_diane_golden.txt"

		cmd := exec.Command("./filter_to_diane.awk", "-v", "var_num=CF00012", "filter_to_diane_testdata.txt")
		var cmdOutput bytes.Buffer
		cmd.Stdout = &cmdOutput
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			t.Fatal(err)
		}

		expectedOutput := createfilter.DiffWithGoldenFile(outGoldenFile, *updateGoldenFile, cmdOutput)
		assert.Equal(t, string(expectedOutput), cmdOutput.String())
	})
}

func TestFilterToDianeXls(t *testing.T) {
	t.Run("filter_to_diane to MS Excel golden file", func(t *testing.T) {

		cmd := exec.Command("ssconvert", "filter_to_diane_golden.txt", "filter_to_diane_golden_tmp.xlsx")
		err := cmd.Run()
		if err != nil {
			t.Fatal(err)
		}

		t.Cleanup(func() {
			_ = os.Remove("filter_to_diane_golden_tmp.xlsx")
		})

		file, err := os.Open("filter_to_diane_golden_tmp.xlsx")
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()
		var cmdOutput bytes.Buffer
		_, err = cmdOutput.ReadFrom(file)
		if err != nil {
			t.Fatal(err)
		}

		cmd = exec.Command("ssconvert", "filter_to_diane_golden_tmp.xlsx", "filter_to_diane_golden_tmp_back.csv")
		err = cmd.Run()
		if err != nil {
			t.Fatal(err)
		}

		t.Cleanup(func() {
			_ = os.Remove("filter_to_diane_golden_tmp_back.csv")
		})

		cmd = exec.Command("grep", "012345678", "filter_to_diane_golden_tmp_back.csv")
		err = cmd.Run()
		if err != nil {
			t.Fatal("The resulting excel file should have kept the leading 0")
		}
	})
}
