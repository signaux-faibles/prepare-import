// Package tools contient du code utilitaire
package tools

import (
	"context"
	"log"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"prepare-import/core"
)

var db *mongo.Database

// SaveInMongo sauve l'objet `toSave` dans une base mongo
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
