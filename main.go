package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/signaux-faibles/prepare-import/common"
)

// Implementation of the prepare-import command.
func main() {
	var path = flag.String("path", ".", "Chemin d'accès aux fichiers données")
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
	validBatchKey, err := common.NewBatchKey(*batchKey)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error()+"\n\nUsage:")
		flag.PrintDefaults()
		os.Exit(1)
	}
	validDateFinEffectif, err := time.Parse("2006-01-02", *dateFinEffectif)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error()+"\n\nUsage:")
		flag.PrintDefaults()
		os.Exit(1)
	}
	adminObject, err := PrepareImport(*path, validBatchKey, dateFinEffectifType(validDateFinEffectif))
	if _, ok := err.(UnsupportedFilesError); ok {
		fmt.Fprintln(os.Stderr, err.Error())
	} else if err != nil {
		log.Fatal(err) // will print in the error output stream and exit
	}
	json, err := json.MarshalIndent(adminObject, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(json))
}
