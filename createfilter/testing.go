package createfilter

import (
	"bytes"
	"io/ioutil"
	"log"
)

// DiffWithGoldenFile compares the output of the execution of a test with the expected output.
func DiffWithGoldenFile(filename string, updateGoldenFile bool, cmdOutput bytes.Buffer) []byte {

	if updateGoldenFile {
		ioutil.WriteFile(filename, cmdOutput.Bytes(), 0644)
	}
	expected, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return expected
}
