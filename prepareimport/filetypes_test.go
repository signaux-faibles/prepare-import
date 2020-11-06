package prepareimport

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractFileTypeFromFilename(t *testing.T) {

	// inspired by https://github.com/golang/go/wiki/TableDrivenTests
	cases := []struct {
		name     string
		category ValidFileType
	}{
		// guessed from urssaf files found on stockage/goub server
		{"Sigfaible_debits.csv", DEBIT},
		{"Sigfaible_cotisdues.csv", COTISATION},
		{"Sigfaible_pcoll.csv", PROCOL},
		{"Sigfaible_etablissement_utf8.csv", ADMIN_URSSAF},
		{"Sigfaible_effectif_siret.csv", EFFECTIF},
		{"Sigfaible_effectif_siren.csv", EFFECTIF_ENT},
		{"Sigfaible_delais.csv", DELAI},
		{"Sigfaible_ccsf.csv", CCSF},

		// guessed from dgefp files
		{"act_partielle_conso_depuis2014_FRANCE.csv", APCONSO},
		{"act_partielle_ddes_depuis2015_FRANCE.csv", APDEMANDE},

		// others
		{"Diane_Export_4.txt", DIANE},
		{"Sigfaibles_debits.csv", DEBIT},
		{"Sigfaibles_debits2.csv", DEBIT},
		{"diane_req_2002.csv", DIANE},
		{"diane_req_dom_2002.csv", DIANE},
		{"effectif_dom.csv", EFFECTIF},
		{"filter_siren_2002.csv", FILTER},
		{"sireneUL.csv", SIRENE_UL},
		{"StockEtablissement_utf8_geo.csv", SIRENE},
	}
	for _, testCase := range cases {
		t.Run("should return "+string(testCase.category)+" for file "+testCase.name, func(t *testing.T) {
			got := ExtractFileTypeFromFilename(testCase.name)
			assert.Equal(t, testCase.category, got)
		})
	}
}

func MakeMetadata(metadataFields MetadataProperty) UploadedFileMeta {
	return UploadedFileMeta{MetaData: metadataFields}
}
