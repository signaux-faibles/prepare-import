package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/signaux-faibles/prepare-import/prepareimport"
)

// Implementation of the prepare-import command.
func main() {
	var path = flag.String("path", ".", "Chemin d'accès au répertoire des batches")
	var batchKey = flag.String(
		"batch",
		"",
		"Clé du batch à importer au format AAMM (année + mois + suffixe optionnel)\n"+
			"Exemple: 1802_1",
	)
	var dateFinEffectif = flag.String(
		"date-fin-effectif",
		"",
		"Date de fin des données \"effectif\" fournies, au format AAAA-MM-JJ (année + mois + jour)\n"+
			"Exemple: 2014-01-01",
	)
	flag.Parse()
	validBatchKey, err := prepareimport.NewBatchKey(*batchKey)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error()+"\n\nUsage:")
		flag.PrintDefaults()
		os.Exit(1)
	}
	adminObject, err := prepareimport.PrepareImport(*path, validBatchKey, *dateFinEffectif)
	if _, ok := err.(prepareimport.UnsupportedFilesError); ok {
		fmt.Fprintln(os.Stderr, err.Error())
	} else if err != nil {
		log.Fatal(err) // will print in the error output stream and exit
	}
	json, err := json.MarshalIndent(adminObject, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(json))
	println("Caution: please make sure that files listed in complete_types were correctly recognized as complete.")
}
