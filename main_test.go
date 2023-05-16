package main

import (
	"github.com/signaux-faibles/prepare-import/createfilter"
	"github.com/signaux-faibles/prepare-import/prepareimport"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var outGoldenFile = "end_to_end_golden.txt"

func Test_prepare(t *testing.T) {
	effectifData, err := os.ReadFile("./createfilter/test_data.csv")
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		batch       string
		finEffectif string
	}
	type want struct {
		adminObject string
		error       string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			"test avec tous les bons paramètres",
			args{"1802", "2018-01-01"},
			want{createfilter.ReadGoldenFile(outGoldenFile), prepareimport.UnsupportedFilesError{}.Error()},
		},
		{
			"test avec un mauvais paramètre batch",
			args{"180", "2018-01-01"},
			want{adminObject: "{}", error: "la clé du batch doit respecter le format requis AAMM"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			buildedBatchKey, _ := prepareimport.NewBatchKey(tt.args.batch)
			parentDir := prepareimport.CreateTempFilesWithContent(t, buildedBatchKey, map[string][]byte{
				"sigfaibles_effectif_siret.csv": effectifData,
				"sigfaibles_debits.csv":         {},
				"abcdef":                        {},
				"abcdef.info":                   []byte(`{ "MetaData": { "filename": "FICHIER_SF_2020_02.csv", "goup-path": "bdf" } }`),
				"unsupported.csv":               {},
				"E_202011095813_Retro-Paydex_20201207.csv": {},
				"083fe617e80f2e30a21598d38a854bc6":         {},
				"083fe617e80f2e30a21598d38a854bc6.info":    []byte(`{ "MetaData": { "filename": "sigfaible_pcoll.csv.gz", "goup-path": "" }, "Size": 1646193 }`),
			})
			object, err2 := prepare(parentDir, tt.args.batch, tt.args.finEffectif)
			assert.ErrorContains(t, err2, tt.want.error)
			assert.Equal(t, tt.want.adminObject, object.ToJSON())
		})
	}
}
