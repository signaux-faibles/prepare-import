// $ go test # to run the tests

package main

import (
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Prepare import should return json object.
func TestPrepareImport(t *testing.T) {
	t.Run("Should return a json with one file", func(t *testing.T) {
		res, _ := PrepareImport([]string{"Sigfaibles_debits.csv"})
		expected := FileProperty{"debit": []string{"Sigfaibles_debits.csv"}}
		assert.Equal(t, expected, res)
		//"{\"files\": {\"debit\": [\"Sigfaibles_debits.csv\"]}}"
	})
	t.Run("Should return an empty json when there is no file", func(t *testing.T) {
		res, _ := PrepareImport([]string{})
		assert.Equal(t, FileProperty{}, res)
	})
}

// To make slices of strings comparable.
func serializeSlice(stringsSlice []string) string {
	stringsClone := append(stringsSlice[:0:0], stringsSlice...)
	sort.Strings(stringsClone)
	return strings.Join(stringsClone, ",")
}

// This function can be used to reduce duplication of assertions,
// while explaining why a failing tests did fail.
// (see `assertCorrectMessage` from https://github.com/quii/learn-go-with-tests/blob/master/hello-world.md#hello-world-again)
func isEqual(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q instead of %q", got, want)
	}
}

func isEqualSlice(t *testing.T, got, want []string) {
	t.Helper()
	if serializeSlice(got) != serializeSlice(want) {
		t.Errorf("got %q instead of %q", got, want)
	}
}

func TestPopulateFilesProperty(t *testing.T) {

	// t.Run() is used to define sub-tests. (see https://golang.org/pkg/testing/#hdr-Subtests_and_Sub_benchmarks)

	t.Run("PopulateFilesProperty should contain effectif file in \"effectif\" property", func(t *testing.T) {
		filesProperty := PopulateFilesProperty([]string{"Sigfaibles_effectif_siret.csv"})
		// isEqualSlice(t, filesProperty["effectif"], []string{"Sigfaibles_effectif_siret.csv"})
    assert.Equal(t,  []string{"Sigfaibles_effectif_siret.csv"}, filesProperty["effectif"])
	})

	t.Run("PopulateFilesProperty should contain one debit file in \"debit\" property", func(t *testing.T) {
		filesProperty := PopulateFilesProperty([]string{"Sigfaibles_debits.csv"})
		assert.Equal(t,  []string{"Sigfaibles_debits.csv"}, filesProperty["debit"])
	})

	t.Run("PopulateFilesProperty should contain both debits files in \"debit\" property", func(t *testing.T) {
		filesProperty := PopulateFilesProperty([]string{"Sigfaibles_debits.csv", "Sigfaibles_debits2.csv"})
		assert.Equal(t,  []string{"Sigfaibles_debits.csv", "Sigfaibles_debits2.csv"}, filesProperty["debit"])
	})
}

func TestGetFileType(t *testing.T) {
	res1, _ := GetFileType("Sigfaibles_effectif_siret.csv")
	if res1 != "effectif" {
		t.Error("GetFileType should return \"effectif\" for \"Sigfaibles_effectif_siret.csv\"")
	}
	res2, _ := GetFileType("Sigfaibles_debits.csv")
	if res2 != "debit" {
		t.Error("GetFileType should return \"debit\" for \"Sigfaibles_debits.csv\"")
	}

	res3, _ := GetFileType("Sigfaibles_debits2.csv")
	if res3 != "debit" {
		t.Error("GetFileType should return \"debit\" for \"Sigfaibles_debits2.csv\"")
	}
}
