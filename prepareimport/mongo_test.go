package prepareimport

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"

	"strconv"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

var mongoURL string
var databaseName string

func Test_SaveInMongo(t *testing.T) {
	ass := assert.New(t)
	batchkeyValue := strconv.Itoa(fake.IntBetween(1000, 9999))
	expectedKey, _ := NewBatchKey(batchkeyValue)
	expectedTypes := []ValidFileType{"effectif", "effectif_ent", "sirene", "sirene_ul", "delai", "procol", "debit", "cotisation"}
	expectedFiles := map[ValidFileType][]string{
		"delai":        {"gzip:/2307/sigfaible_delais.csv.gz"},
		"admin_urssaf": {"gzip:/2307/sigfaible_etablissement_utf8.csv.gz"},
		"procol":       {"gzip:/2307/sigfaible_pcoll.csv.gz"},
		"effectif_ent": {"gzip:/2307/sigfaible_effectif_siren.csv.gz"},
		"effectif":     {"gzip:/2307/sigfaible_effectif_siret.csv.gz"},
		"filter":       {"/2307/filter_siren_2307.csv"},
		"sirene":       {"/2307/StockEtablissement_utf8_geo.csv"},
		"debit":        {"gzip:/2307/sigfaible_debits.csv.gz"},
		"sirene_ul":    {"/2307/sireneUL.csv"},
		"cotisation":   {"gzip:/2307/sigfaible_cotisdues.csv.gz"},
	}
	now := time.Now()
	expectedDateFinEffectif := time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC)
	toSave := AdminObject{
		ID: IDProperty{
			Key:  expectedKey,
			Type: "batch",
		},
		CompleteTypes: expectedTypes,
		Files:         expectedFiles,
		Param: ParamProperty{
			DateDebut:       now.AddDate(-4, 0, 0),
			DateFin:         now,
			DateFinEffectif: expectedDateFinEffectif,
		},
	}

	id, err := SaveInMongo(context.Background(), toSave, mongoURL, databaseName)
	ass.NoError(err)
	key, err := findKey(id)
	ass.NoError(err)

	object, err := fetchMongoAdminObject(key)
	ass.NoError(err)
	// complete types

	ass.ElementsMatch(object.CompleteTypes, expectedTypes)

	// expectedFiles
	ass.NotNil(object.Files)
	ass.Exactly(object.Files, expectedFiles)

	// param
	ass.NotNil(object.Param)
	ass.NotNil(object.Param.DateFinEffectif)
	ass.Equal(expectedDateFinEffectif, object.Param.DateFinEffectif)
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

func fetchMongoAdminObject(batchkey string) (AdminObject, error) {
	key := bson.M{"_id.key": batchkey}
	found := db.Collection("Admin").FindOne(context.Background(), key)
	var result AdminObject
	err := found.Decode(&result)
	if err != nil {
		return AdminObject{}, errors.WithMessage(err, "erreur lors du décodage")
	}
	return result, nil
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
