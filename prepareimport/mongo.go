package prepareimport

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

// SaveInMongo sauve l'objet `toSave` dans une base mongo
func SaveInMongo(ctx context.Context, toSave AdminObject, mongoURL, databaseName string) (interface{}, error) {
	if db == nil {
		err := connectDB(ctx, mongoURL, databaseName)
		if err != nil {
			return nil, errors.Wrap(err, "Erreur pendant la connexion à Mongo")
		}
	}
	saved, err := db.Collection("Admin").InsertOne(ctx, toSave)
	if err != nil {
		return nil, errors.Wrap(err, "Erreur pendant l'insertion du document")
	}
	slog.Info("import sauvé dans la collection `Admin`", slog.Any("id", saved))
	return saved.InsertedID, nil
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

func findKey(current interface{}) (string, error) {
	insertedID, ok := current.(primitive.D)
	if !ok {
		return "", errors.New("erreur pendant la recherche de la clé : mauvais type")
	}
	marshal, _ := bson.Marshal(insertedID)
	var result = primitive.M{}
	_ = bson.Unmarshal(marshal, &result)
	value, found := result["key"]
	if !found {
		return "", errors.New("erreur pendant la recherche de la clé: mauvais id")
	}
	key := fmt.Sprint(value)
	return key, nil
}
