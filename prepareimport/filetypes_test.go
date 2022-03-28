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
		{"Sigfaible_debits.csv", debit},
		{"Sigfaible_cotisdues.csv", cotisation},
		{"Sigfaible_pcoll.csv", procol},
		{"Sigfaible_etablissement_utf8.csv", adminUrssaf},
		{"Sigfaible_effectif_siret.csv", effectif},
		{"Sigfaible_effectif_siren.csv", effectifEnt},
		{"Sigfaible_delais.csv", delai},
		{"Sigfaible_ccsf.csv", ccsf},

		// compressed version of urssaf files
		{"Sigfaible_ccsf.csv.gz", ccsf},
		{"Sigfaible_etablissement_utf8.csv.gz", adminUrssaf}, // sfdata parser name: "comptes"
		{"Sigfaible_cotisdues.csv.gz", cotisation},
		{"Sigfaible_debits.csv.gz", debit},
		{"Sigfaible_delais.csv.gz", delai},
		{"Sigfaible_effectif_siren.csv.gz", effectifEnt},
		{"Sigfaible_effectif_siret.csv.gz", effectif},
		{"Sigfaible_pcoll.csv.gz", procol},

		// guessed from dgefp files
		{"consommation_ap.csv", apconso},
		{"demande_ap.csv", apdemande},

		// others
		{"Diane_Export_4.txt", diane},
		{"Sigfaibles_debits.csv", debit},
		{"Sigfaibles_debits2.csv", debit},
		{"diane_req_2002.csv", diane},
		{"diane_req_dom_2002.csv", diane},
		{"effectif_dom.csv", effectif},
		{"filter_siren_2002.csv", filter},
		{"sireneUL.csv", sireneUl},
		{"StockEtablissement_utf8_geo.csv", sirene},
		{"StockEtablissement_utf8_geo.csv", sirene},
		{"E_202011095813_Retro-Paydex_20201207.csv", paydex},
		{"E_202011095813_Identite_20201207.csv", ""}, // not paydex
		{"Ellisphère-Tête de groupe-FinalV2-2015.xlsx", ellisphere},
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
