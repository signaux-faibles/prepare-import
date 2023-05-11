package createfilter

import (
	"bytes"
	"log"
	"os"
)

// DiffWithGoldenFile compares the output of the execution of a test with the expected output.
func DiffWithGoldenFile(filename string, updateGoldenFile bool, cmdOutput bytes.Buffer) []byte {

	if updateGoldenFile {
		os.WriteFile(filename, cmdOutput.Bytes(), 0644)
	}
	expected, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return expected
}
