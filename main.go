// $ go build # to compile
// $ ./prepare-import # to run

package main

import (
	"errors"
  "io/ioutil"
)

func main() {
}

func ReadFilenames(path string) ([]string, error) {
  var files []string
  fileInfo, err := ioutil.ReadDir(path)
  if err != nil {
    return files, err
  }
  for _, file := range fileInfo {
    files = append(files, file.Name())
  }
  return files, nil
}

func PrepareImport(filenames []string) (FileProperty, error) {
  // func Valid(data []byte) bool
  fileProperty := PopulateFilesProperty(filenames)

  return fileProperty, nil
}

type FileProperty map[string][]string

func PopulateFilesProperty(filenames []string) FileProperty {
  fileProperty := FileProperty{
    // "effectif": []string{"coucou"},
    // "debit":    []string{},
  }
  for _, filename := range filenames {
    filetype, _ := GetFileType(filename)
    if _, exists := fileProperty[filetype]; !exists {
      fileProperty[filetype] = []string{}
    }
    fileProperty[filetype] = append(fileProperty[filetype], filename)
  }
  return fileProperty
}

func GetFileType(filename string) (string, error) {
  switch filename {
  case "Sigfaibles_effectif_siret.csv":
    return "effectif", nil
  case "Sigfaibles_debits.csv":
    return "debit", nil
  case "Sigfaibles_debits2.csv":
    return "debit", nil
  default:
    return "", errors.New("Unrecognized type for " + filename)
  }
}
