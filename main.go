package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/pkg/errors"

	"prepare-import/prepareimport"
)

var loglevel *slog.LevelVar

// Implementation of the prepare-import command.
func main() {
	InitLogger(slog.LevelInfo)
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
	adminObject, err := prepare(*path, *batchKey, *dateFinEffectif)
	if err != nil {
		slog.Error("erreur pendant la génération de l'objet", slog.Any("error", err))
		panic(err)
	}
	slog.Info("objet généré", slog.Any("adminObject", adminObject))
	saveAdminObject(adminObject, *mongoURL, *databaseName)
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
		return prepareimport.AdminObject{}, errors.Wrap(err, "erreur inattendue pendant la préparation de l'import")
	}
	return adminObject, nil
}

func saveAdminObject(toSave prepareimport.AdminObject, mongoURL string, databaseName string) {
	if len(mongoURL) <= 0 || len(databaseName) <= 0 {
		_, _ = fmt.Fprintln(os.Stderr, "ATTENTION : le résultat ne sera pas sauvegardé en base car au moins un des paramètres nécessaires n'a pas été spécifié")
		return
	}
	_, err := prepareimport.SaveInMongo(context.Background(), toSave, mongoURL, databaseName)
	if err != nil {
		slog.Error("Erreur inattendue pendant la sauvegarde de l'import : ", slog.Any("error", err))
		panic(err)
	}
}

func InitLogger(level slog.Level) {
	loglevel = new(slog.LevelVar)
	loglevel.Set(level)
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: loglevel,
	})

	appLogger := slog.New(handler)
	slog.SetDefault(appLogger)
}
