package createfilter

import (
	"bytes"
	"encoding/csv"
	"flag"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var outGoldenFile = "test_golden.txt"
var errGoldenFile = "test_golden_err.txt"

var updateGoldenFile = flag.Bool("update", false, "Update the expected test values in golden file")

func TestCreateFilter(t *testing.T) {
	t.Run("create_filter golden file", func(t *testing.T) {

		var cmdOutput bytes.Buffer
		var cmdError bytes.Buffer = *bytes.NewBufferString("") // default: no error
		err := CreateFilter(&cmdOutput, "test_data.csv", DefaultNbMois, DefaultMinEffectif, DefaultNbIgnoredCols)
		if err != nil {
			cmdError = *bytes.NewBufferString(err.Error())
		}

		expectedOutput := DiffWithGoldenFile(outGoldenFile, *updateGoldenFile, cmdOutput)
		expectedError := DiffWithGoldenFile(errGoldenFile, *updateGoldenFile, cmdError)

		assert.Equal(t, string(expectedOutput), cmdOutput.String())
		assert.Equal(t, string(expectedError), cmdError.String())
	})
}

func TestIsInsidePerimeter(t *testing.T) {
	nbMois := 3
	minEffectif := 10
	testCases := []struct {
		input    []string
		expected bool
	}{
		{[]string{"10", "9", "4", "7", "5"}, false},
		{[]string{"10", "20", "4", "7", "5"}, false},
		{[]string{"10", "9", "12", "7", "5"}, true},
		{[]string{"10", "9", "12", "", ""}, true},
		{[]string{"10", "9", "5", "", ""}, false},
		{[]string{"10", "9", "", "", ""}, false},
	}

	for i, tc := range testCases {
		t.Run("Test case "+strconv.Itoa(i), func(t *testing.T) {
			shouldKeep := isInsidePerimeter(tc.input, nbMois, minEffectif)
			assert.Equal(t, tc.expected, shouldKeep)
		})
	}
}

func TestGuessLastNonMissing(t *testing.T) {
	testCases := []struct {
		inputCsv string
		expected int
	}{
		{"1,", 1},
		{",1", 0},
		{"1,1", 0},
		{",", 2},
		{",\n,1", 0},
		{"1,\n,", 1},
		{"1,\n1,", 1},
	}

	for i, tc := range testCases {
		t.Run("Test case without ignored "+strconv.Itoa(i), func(t *testing.T) {
			reader := csv.NewReader(strings.NewReader(tc.inputCsv))
			lastNonMissing := guessLastNMissingFromReader(reader, 0)
			assert.Equal(t, tc.expected, lastNonMissing)
		})
	}

	testCasesIgnore := []struct {
		inputCsv string
		expected int
	}{
		{"1,,1", 1},
		{",1,1", 0},
		{"1,1,1", 0},
		{",,1", 2},
		{",,1\n,1,1", 0},
		{"1,,1\n,,1", 1},
		{"1,,1\n1,,1", 1},
	}

	for i, tc := range testCasesIgnore {
		t.Run("Test case without ignored "+strconv.Itoa(i), func(t *testing.T) {
			reader := csv.NewReader(strings.NewReader(tc.inputCsv))
			lastNonMissing := guessLastNMissingFromReader(reader, 1)
			assert.Equal(t, tc.expected, lastNonMissing)
		})
	}
}
