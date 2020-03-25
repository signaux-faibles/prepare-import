// $ go build # to compile
// $ ./prepare-import # to run

package main

import ("errors")

func main(){
}

func PrepareImport() (string, error) {
	// func Valid(data []byte) bool
	return "{}", nil
}

type FileProperty map[string][]string

func PopulateFilesProperty(filenames []string) FileProperty {
  fileProperty := FileProperty{}
  for _, filename := range filenames {
    filetype, _ := GetFileType(filename)
    fileProperty[filetype] = []string{}
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
