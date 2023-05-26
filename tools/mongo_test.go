package tools

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/jaswdr/faker"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"

	"prepare-import/core"
)

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

	fmt.Println("On peut lancer les tests")
	code := m.Run()
	//killContainer(mongodb)
	// You can't defer this because os.Exit doesn't care for defer

	os.Exit(code)
}

func Test_SaveInMongo(t *testing.T) {
	ass := assert.New(t)
	firstID := fake.Lorem().Word()
	secondID := fake.Lorem().Word()
	t.Run("on sauve un object avec le premier `_id`", func(t *testing.T) {
		toSave := core.AdminObject{"_id": firstID, fake.Beer().Name(): fake.Beer().Hop()}
		err := SaveInMongo(context.Background(), toSave, mongoURL, databaseName)
		ass.NoError(err)
	})
	t.Run("on sauve une objet fois avec un autre `_id`", func(t *testing.T) {
		toSave := core.AdminObject{"_id": secondID, fake.App().Name(): fake.App().Version()}
		err := SaveInMongo(context.Background(), toSave, mongoURL, databaseName)
		ass.NoError(err)
	})
	t.Run("on sauve un autre objet avec le premier `_id` une seconde fois", func(t *testing.T) {
		toSave := core.AdminObject{"_id": firstID, fake.Internet().TLD(): fake.Internet().Domain()}
		err := SaveInMongo(context.Background(), toSave, mongoURL, databaseName)
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
			Tag:        "6.0",
			Env: []string{
				// username and password for mongodb superuser
				"MONGO_INITDB_ROOT_USERNAME=root",
				"MONGO_INITDB_ROOT_PASSWORD=password",
			},
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
