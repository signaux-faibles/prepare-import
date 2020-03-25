// $ go build # to compile
// $ ./prepare-import # to run

package main

func main(){
}

func PrepareImport() (string, error) {
	// func Valid(data []byte) bool
	return "{}", nil
}

func GetFileType(filename string) (string) {
	switch filename {
	case "darwin":
		fmt.Println("OS X.")
	case "linux":
		fmt.Println("Linux.")
	default:
		// freebsd, openbsd,
		// plan9, windows...
		fmt.Printf("%s.\n", os)
	}
	return "effectif"
}
