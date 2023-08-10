package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"

	"prepare-import/core"
	"prepare-import/prepareimport"
	"prepare-import/tools"
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
	var mongoURL = flag.String("mongoURL", "", "Url de connexion à la base Mongo\n"+
		"Exemple: mongodb://username:password@ip:port")
	var databaseName = flag.String("databaseName", "", "Nom de la base de données Mongo")

	flag.Parse()
	adminObject, _ := prepare(*path, *batchKey, *dateFinEffectif)
	saveAdminObject(adminObject, *mongoURL, *databaseName)
	println("Caution: please make sure that files listed in complete_types were correctly recognized as complete.")
}

func prepare(path, batchKey, dateFinEffectif string) (core.AdminObject, error) {
	validBatchKey, err := prepareimport.NewBatchKey(batchKey)
	if err != nil {
		return core.AdminObject{}, errors.Wrap(err, "erreur lors de la création de la clé de batch")
	}
	adminObject, err := prepareimport.PrepareImport(path, validBatchKey, dateFinEffectif)
	if _, ok := err.(prepareimport.UnsupportedFilesError); ok {
		return adminObject, err
	} else if err != nil {
		return core.AdminObject{}, errors.Wrap(err, "erreur inattendue pendant la préparation de l'import : ")
	}
	return adminObject, nil
}

func saveAdminObject(toSave core.AdminObject, mongoURL string, databaseName string) {
	println(toSave.ToJSON())
	if len(mongoURL) <= 0 || len(databaseName) <= 0 {
		_, _ = fmt.Fprintln(os.Stderr, "ATTENTION : le résultat ne sera pas sauvegardé en base car au moins un des paramètres nécessaires n'a pas été spécifié")
		println(toSave.ToJSON())
		return
	}
	err := tools.SaveInMongo(context.Background(), toSave, mongoURL, databaseName)
	if err != nil {
		log.Fatal("Erreur inattendue pendant la sauvegarde de l'import : ", err)
	}
}
