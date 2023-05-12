package main

import (
	"bytes"
	"flag"
	"github.com/jaswdr/faker"
	"github.com/signaux-faibles/prepare-import/createfilter"
	"github.com/signaux-faibles/prepare-import/prepareimport"
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"strconv"
	"testing"
)

var outGoldenFile = "end_to_end_golden.txt"
var errGoldenFile = "end_to_end_golden_err.txt"

var updateGoldenFile = flag.Bool("update", false, "Update the expected test values in golden file")

var fake faker.Faker

func init() {
	fake = faker.New()
}

func Test_Main(t *testing.T) {

	t.Run("prepare-import golden file", func(t *testing.T) {
		t.Log("ATTENTION: ce test utilise l'exécutable compilé et non les sources.")
		t.Log("Il faut donc builder pour etre sur de tester la bonne version")
		batch := "1802"

		batchKey, _ := prepareimport.NewBatchKey(batch)

		effectifData, err := os.ReadFile("./createfilter/test_data.csv")
		if err != nil {
			t.Fatal(err)
		}
		parentDir := prepareimport.CreateTempFilesWithContent(t, batchKey, map[string][]byte{
			"sigfaibles_effectif_siret.csv": effectifData,
			"sigfaibles_debits.csv":         {},
			"abcdef":                        {},
			"abcdef.info":                   []byte(`{ "MetaData": { "filename": "FICHIER_SF_2020_02.csv", "goup-path": "bdf" } }`),
			"unsupported.csv":               {},
			"E_202011095813_Retro-Paydex_20201207.csv": {},
			"083fe617e80f2e30a21598d38a854bc6":         {},
			"083fe617e80f2e30a21598d38a854bc6.info":    []byte(`{ "MetaData": { "filename": "sigfaible_pcoll.csv.gz", "goup-path": "" }, "Size": 1646193 }`),
		})

		//url := randomMongoURL()
		//databaseName := fake.Lorem().Word()

		cmds := []*exec.Cmd{
			exec.Command(
				"./prepare-import",
				"-path", parentDir,
				"-batch", batch,
				"-date-fin-effectif", "2018-01-01",
				//"-mongoURL", url,
				//"-databaseName", databaseName,
			), // paramètres valides
			exec.Command("./prepare-import", "-path", parentDir, "-batch", "180"), // nom de batch invalide
		}
		var cmdOutput bytes.Buffer
		var cmdError bytes.Buffer
		for _, cmd := range cmds {
			cmd.Stdout = &cmdOutput
			cmd.Stderr = &cmdError
			_ = cmd.Run()
		}

		expectedOutput := createfilter.DiffWithGoldenFile(outGoldenFile, *updateGoldenFile, cmdOutput)
		expectedError := createfilter.DiffWithGoldenFile(errGoldenFile, *updateGoldenFile, cmdError)

		assert.Equal(t, string(expectedOutput), cmdOutput.String())
		assert.Equal(t, string(expectedError), cmdError.String())
	})
}

func randomMongoURL() string {
	return "mongodb://" + fake.Internet().User() + ":" + fake.Internet().Password() + "@" + fake.Internet().Ipv4() + ":" + strconv.Itoa(fake.IntBetween(4000, 8000))
}
