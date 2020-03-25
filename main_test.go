// $ go test # to run the tests

package main

import ("testing")

// Prepare import should return json object.
func TestPrepareImport(t *testing.T){
  res, _ := PrepareImport()
  if (res != "{}") {
    t.Error("Test failed: invalid json")
  }
}

func TestGetFileType(t *testing.T){
	res1 := GetFileType("Sigfaibles_effectif_siret.csv")
	if (res1 != "effectif") {
		t.Error("GetFileType should return \"effectif\" for \"Sigfaibles_effectif_siret.csv\"")
	}
	res2 := GetFileType("Sigfaibles_debits.csv")
  if (res2 != "debit") {
    t.Error("GetFileType should return \"debit\" for \"Sigfaibles_debits.csv\"")
  }
}
