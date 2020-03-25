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
