package tools

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/jaswdr/faker"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"prepare-import/core"
)

//go:embed adminObject.json
var adminObjectJson string

////go:embed file.txt
//var adminObjectBson []byte

var mongoURL string
var databaseName string
var fake faker.Faker

func init() {
	fake = faker.New()
}
func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		panic("ne peut démarrer mongodb")
	}

	mongoURL, databaseName = startMongoDB(pool)

	fmt.Println("On peut lancer les tests, mongo:", mongoURL, " / dbname:", databaseName)
	code := m.Run()
	//killContainer(mongodb)
	// You can't defer this because os.Exit doesn't care for defer

	os.Exit(code)
}

func Test_SaveInMongo_fromJSON(t *testing.T) {
	ass := assert.New(t)
	toSave := core.FromJSON(adminObjectJson)
	id, err := SaveInMongo(context.Background(), toSave, mongoURL, databaseName)
	ass.NoError(err)
	key, err := findKey(id)
	ass.NoError(err)
	object, err := fetchMongoAdminObject(key)
	ass.NoError(err)
	// complete types
	completeTypes, good := object["complete_types"].(primitive.A)
	ass.True(good)
	ass.Len(completeTypes, 8)
	ass.Contains(completeTypes, "effectif")
	ass.ElementsMatch(completeTypes, []string{"effectif", "effectif_ent", "sirene", "sirene_ul", "cotisation", "delai", "procol", "debit"})

	// files
	ass.IsType(primitive.M{}, object["files"])
	files := object["files"].(primitive.M)
	ass.NotNil(files)
	ass.Len(files, 10)
	//on en teste 1 seul
	ass.ElementsMatch(files["effectif_ent"], []string{"gzip:/2307/sigfaible_effectif_siren.csv.gz"})

	// param
	param, good := object["param"].(primitive.M)
	ass.True(good)
	ass.NotNil(param)
	ass.Len(param, 3)
	ass.NotNil(param["date_fin_effectif"])
	// TODO asserter sur le type date
	//spew.Dump(object)
}

func Test_SaveInMongo_loadAndWriteFile(t *testing.T) {
	key := bson.M{"_id.key": "2307"}
	found := db.Collection("Admin").FindOne(context.Background(), key)
	var result interface{}
	_ = found.Decode(&result)
	data, _ := bson.Marshal(result)
	if err := os.WriteFile("file.txt", data, 0666); err != nil {
		log.Fatal(err)
	}

}

func Test_SaveInMongo_fromBSON(t *testing.T) {
	//ass := assert.New(t)
	//toSave := core.FromBSON(adminObjectBson)
	//id, err := SaveInMongo(context.Background(), toSave, mongoURL, databaseName)
	//ass.NoError(err)
	//key, err := findKey(id)
	//ass.NoError(err)
	//object, err := fetchMongoAdminObject(key)
	//ass.NoError(err)
	//// complete types
	//completeTypes, good := object["complete_types"].(primitive.A)
	//ass.True(good)
	//ass.Len(completeTypes, 8)
	//ass.Contains(completeTypes, "effectif")
	//ass.ElementsMatch(completeTypes, []string{"effectif", "effectif_ent", "sirene", "sirene_ul", "cotisation", "delai", "procol", "debit"})
	//
	//// files
	//ass.IsType(primitive.M{}, object["files"])
	//files := object["files"].(primitive.M)
	//ass.NotNil(files)
	//ass.Len(files, 10)
	////on en teste 1 seul
	//ass.ElementsMatch(files["effectif_ent"], []string{"gzip:/2307/sigfaible_effectif_siren.csv.gz"})
	//
	//// param
	//param, good := object["param"].(primitive.M)
	//ass.True(good)
	//ass.NotNil(param)
	//ass.Len(param, 3)
	//ass.NotNil(param["date_fin_effectif"])
	//// TODO asserter sur le type date
	////spew.Dump(object)
}

func fetchMongoAdminObject(batchkey string) (core.AdminObject, error) {
	key := bson.M{"_id.key": batchkey}
	found := db.Collection("Admin").FindOne(context.Background(), key)
	var result interface{}
	err := found.Decode(&result)
	if err != nil {
		return nil, errors.WithMessage(err, "erreur lors du décodage")
	}
	data, _ := bson.Marshal(result)
	mappe := bson.M{}
	_ = bson.Unmarshal(data, mappe)
	return core.AdminObject(mappe), nil
}

func Test_SaveInMongo(t *testing.T) {
	ass := assert.New(t)
	firstID := fake.Lorem().Word()
	secondID := fake.Lorem().Word()
	t.Run("on sauve un object avec le premier `_id`", func(t *testing.T) {
		toSave := core.AdminObject{"_id": firstID, fake.Beer().Name(): fake.Beer().Hop()}
		_, err := SaveInMongo(context.Background(), toSave, mongoURL, databaseName)
		ass.NoError(err)
	})
	t.Run("on sauve un objet avec un autre `_id`", func(t *testing.T) {
		toSave := core.AdminObject{"_id": secondID, fake.App().Name(): fake.App().Version()}
		_, err := SaveInMongo(context.Background(), toSave, mongoURL, databaseName)
		ass.NoError(err)
	})
	t.Run("on sauve un autre objet avec le premier `_id` une seconde fois", func(t *testing.T) {
		toSave := core.AdminObject{"_id": firstID, fake.Internet().TLD(): fake.Internet().Domain()}
		_, err := SaveInMongo(context.Background(), toSave, mongoURL, databaseName)
		ass.Error(err)
		ass.ErrorContains(err, "E11000 duplicate key error collection")
	})
}

func startMongoDB(pool *dockertest.Pool) (mongoURL, databaseName string) {

	// pulls an image, creates a container based on it and runs it
	mongodbContainerName := "mongodb-ti-" + strconv.Itoa(time.Now().Nanosecond())

	mongodb, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Name:       mongodbContainerName,
			Repository: "mongo",
			Tag:        "4.2",
			Env: []string{
				// username and password for mongodb superuser
				"MONGO_INITDB_ROOT_USERNAME=root",
				"MONGO_INITDB_ROOT_PASSWORD=password",
			},
			//Cmd: []string{
			//	"--wiredTigerJournalCompressor zstd",
			//	//--wiredTigerCollectionBlockCompressor zstd --wiredTigerIndexPrefixCompression true --wiredTigerCacheSizeGB 4",
			//},
			//Mounts: []string{cwd + "/test/resources/:/dump/"},
		},
		func(config *docker.HostConfig) {
			// set AutoRemove to true so that stopped container goes away by itself
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{
				Name: "no",
			}
		},
	)
	if err != nil {
		fmt.Println(err.Error())
		killContainer(mongodb)
	}

	// container stops after 600 seconds
	if err = mongodb.Expire(600); err != nil {
		fmt.Print("Erreur lors de de l'application d'une date d'expiration sur le container mongo : ", err)
		killContainer(mongodb)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		fmt.Println("Mongo n'est pas encore prêt")
		var err error
		databaseName = "test"
		mongoURL = fmt.Sprintf("mongodb://root:password@localhost:%s", mongodb.GetPort("27017/tcp"))
		err = connectDB(context.Background(), mongoURL, databaseName)
		if err != nil {
			return err
		}
		return pingDB()
	}); err != nil {
		fmt.Printf("N'arrive pas à démarrer/restaurer Mongo: %s", err)
	}
	return mongoURL, databaseName
}

func killContainer(resource *dockertest.Resource) {
	if resource == nil {
		return
	}
	if err := resource.Close(); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}
