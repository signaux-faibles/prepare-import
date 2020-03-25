package main

import ("testing")
// Prepare import should return json object.
func TestPrepareImport(t *testing.T){
  res := PrepareImport()

  if (res != "{}") {
    t.Error("Test failed: invalid json")
  }


}
