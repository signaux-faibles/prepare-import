package createfilter

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"flag"
	"strconv"
	"strings"
	"testing"
	"time"

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

func TestOutputPerimeter(t *testing.T) {
	t.Run("le département de l'entreprise n'est pas considéré comme une valeur d'effectif", func(t *testing.T) {
		// setup conditions and expectations
		minEffectif := 10
		nbIgnoredCols := 2 // "base" and "UR_EMET"
		expectedSirens := []string{"22222222222222", "33333333333333"}
		effectifData := strings.Join([]string{
			"compte;siret;rais_soc;ape_ins;dep;eff201011;eff201012;base;UR_EMET",
			"000000000000000000;00000000000000;ENTREPRISE;1234Z;75;4;4;116;075077",   // ❌ 75 ≥ 10, mais ce n'est pas un effectif
			"111111111111111111;11111111111111;ENTREPRISE;1234Z;53;4;4;116;075077",   // ❌ 53 ≥ 10, mais ce n'est pas un effectif
			"222222222222222222;22222222222222;ENTREPRISE;1234Z;92;14;14;116;075077", // ✅ siren retenu car 14 est bien un effectif ≥ 10
			"333333333333333333;33333333333333;ENTREPRISE;1234Z;92;14;14;116;075077", // ✅ siren retenu car 14 est bien un effectif ≥ 10
		}, "\n")
		// run the test
		var output bytes.Buffer
		reader := csv.NewReader(strings.NewReader(effectifData))
		reader.Comma = ';'
		writer := bufio.NewWriter(&output)
		outputPerimeter(reader, writer, DefaultNbMois, minEffectif, nbIgnoredCols)
		writer.Flush()
		// assert
		assert.Equal(t, expectedSirens, strings.Split(output.String(), "\n"))
	})
}

func TestDetectDateFinEffectif(t *testing.T) {
	expectedDate := time.Date(2020, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
	actualDate, err := DetectDateFinEffectif("test_data.csv", DefaultNbIgnoredCols) // => col name: "eff202011"
	if assert.NoError(t, err) {
		assert.Equal(t, expectedDate, actualDate)
	}
}

func TestIsInsidePerimeter(t *testing.T) {
	nbMois := 3 // => number of recent months that will be considered
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
