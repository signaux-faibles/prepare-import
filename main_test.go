// $ go test # to run the tests

package main

import (
	"sort"
	"strings"
	"testing"
)

// Prepare import should return json object.
func TestPrepareImport(t *testing.T) {
	res, _ := PrepareImport()
	if res != "{}" {
		t.Error("Test failed: invalid json")
	}
}

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

func TestPopulateFilesProperty(t *testing.T) {

	/**
	 * t.Run() can be used to define sub-tests. (see https://golang.org/pkg/testing/#hdr-Subtests_and_Sub_benchmarks)
	 *
	 * Here's what is logged by `$ go test` if the first t.Run() block fails:
	 *
	 * --- FAIL: TestPopulateFilesProperty (0.00s)
	 *     --- FAIL: TestPopulateFilesProperty/PopulateFilesProperty_should_contain_effectif_file_in_"effectif"_property (0.00s)
	 *         main_test.go:31: PopulateFilesProperty should contain effectif file in "effectif" property
	 */

	t.Run("PopulateFilesProperty should contain effectif file in \"effectif\" property", func(t *testing.T) {
		filesProperty := PopulateFilesProperty([]string{"Sigfaibles_effectif_siret.csv"})
		if _, ok := filesProperty["effectif"]; !ok {
			t.Error("PopulateFilesProperty should have a \"effectif\" property")
		}
		isEqual(t, filesProperty["effectif"][0], "Sigfaibles_effectif_siret.csv")
	})

	t.Run("PopulateFilesProperty should contain one debit file in \"debit\" property", func(t *testing.T) {
		filesProperty := PopulateFilesProperty([]string{"Sigfaibles_debits.csv"})
		if _, ok := filesProperty["debit"]; !ok {
			t.Error("PopulateFilesProperty should have a \"debit\" property")
		}
		isEqual(t, filesProperty["debit"][0], "Sigfaibles_debits.csv")
	})

	t.Run("PopulateFilesProperty should contain both debits files in \"debit\" property", func(t *testing.T) {
		filesProperty := PopulateFilesProperty([]string{"Sigfaibles_debits.csv", "Sigfaibles_debits2.csv"})
		if _, ok := filesProperty["debit"]; !ok {
			t.Error("PopulateFilesProperty should have a \"debit\" property")
		}
		isEqual(t, filesProperty["debit"][0], "Sigfaibles_debits.csv")
		isEqual(t, filesProperty["debit"][1], "Sigfaibles_debits2.csv")
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
