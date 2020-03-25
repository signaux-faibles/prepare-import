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

func TestPopulateFilesProperty(t *testing.T){
  var filesProperty FileProperty
  filesProperty = PopulateFilesProperty([]string{"Sigfaibles_effectif_siret.csv"})
  if (filesProperty["effectif"][0] != "Sigfaibles_effectif_siret.csv") {
    t.Error("PopulateFilesProperty should contain effectif file in \"effectif\" property")
	}
	
	filesProperty = PopulateFilesProperty([]string{"Sigfaibles_debits.csv"})
	if _, ok := filesProperty["debit"]; !ok {
    t.Error("PopulateFilesProperty should return \"debit\" property")
	}
  if (filesProperty["debit"][0] != "Sigfaibles_debits.csv") {
    t.Error("PopulateFilesProperty should contain debits file in \"debit\" property")
	}

}

func TestGetFileType(t *testing.T){
	res1, _ := GetFileType("Sigfaibles_effectif_siret.csv")
	if (res1 != "effectif") {
		t.Error("GetFileType should return \"effectif\" for \"Sigfaibles_effectif_siret.csv\"")
	}
	res2, _ := GetFileType("Sigfaibles_debits.csv")
  if (res2 != "debit") {
    t.Error("GetFileType should return \"debit\" for \"Sigfaibles_debits.csv\"")
  }

  res3, _ := GetFileType("Sigfaibles_debits2.csv")
  if (res3 != "debit") {
    t.Error("GetFileType should return \"debit\" for \"Sigfaibles_debits2.csv\"")
  }
}
