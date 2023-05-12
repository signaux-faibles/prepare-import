package tools

import (
	"context"
	"github.com/pkg/errors"
	"github.com/signaux-faibles/prepare-import/core"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var db *mongo.Database

func SaveInMongo(ctx context.Context, toSave core.AdminObject, mongoURL, databaseName string) error {
	if db == nil {
		err := connectDB(ctx, mongoURL, databaseName)
		if err != nil {
			return errors.Wrap(err, "Erreur pendant la connexion à Mongo")
		}
	}
	saved, err := db.Collection("Admin").InsertOne(ctx, toSave)
	if err != nil {
		return errors.Wrap(err, "Erreur pendant l'insertion du document")
	}
	log.Printf("import sauvé dans la collection `Admin` avec l'id : %s", saved)
	return nil
}

func pingDB() error {
	return db.Client().Ping(context.Background(), nil)
}

func connectDB(ctx context.Context, uri string, databaseName string) error {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}
	db = client.Database(databaseName)
	return nil
}
