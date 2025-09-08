package main

import (
	"flag"
	"log"
	"path"

	"github.com/pkg/errors"

	"prepare-import/prepareimport"
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
	var configFile = flag.String("configFile", "./config.toml", "Chemin du fichier où est écrit la configuration,"+
		" relatif au chemin d'accès au répertoire des batches \n"+
		"Exemple: ./config.toml")

	flag.Parse()
	adminObject, err := prepare(*path, *batchKey, *dateFinEffectif)
	if err != nil {
		panic(err)
	}
	saveAdminObject(adminObject, *path, *configFile)
	println("Caution: please make sure that files listed in complete_types were correctly recognized as complete.")
}

func prepare(path, batchKey, dateFinEffectif string) (prepareimport.AdminObject, error) {
	validBatchKey, err := prepareimport.NewBatchKey(batchKey)
	if err != nil {
		return prepareimport.AdminObject{}, errors.Wrap(err, "erreur lors de la création de la clé de batch")
	}
	adminObject, err := prepareimport.PrepareImport(path, validBatchKey, dateFinEffectif)
	if _, ok := err.(prepareimport.UnsupportedFilesError); ok {
		return adminObject, err
	} else if err != nil {
		return prepareimport.AdminObject{}, errors.Wrap(err, "erreur inattendue pendant la préparation de l'import : ")
	}
	return adminObject, nil
}

func saveAdminObject(toSave prepareimport.AdminObject, dirPath, configFile string) {
	filePath := path.Join(dirPath, configFile)
	err := prepareimport.SaveToFile(toSave, filePath)

	if err != nil {
		log.Fatal("Erreur inattendue pendant la sauvegarde de l'import : ", err)
	}
}
