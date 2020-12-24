package prepareimport

import (
	"regexp"
	"strings"
)

// ExtractFileTypeFromFilename returns a file type from filename, or empty string for unsupported file names
func ExtractFileTypeFromFilename(filename string) ValidFileType {
	switch {
	case filename == "act_partielle_conso_depuis2014_FRANCE.csv":
		return apconso
	case filename == "act_partielle_ddes_depuis2015_FRANCE.csv":
		return apdemande
	case filename == "Sigfaible_etablissement_utf8.csv":
		return adminUrssaf
	case filename == "Sigfaible_effectif_siren.csv":
		return effectifEnt
	case filename == "Sigfaible_pcoll.csv":
		return procol
	case filename == "Sigfaible_cotisdues.csv":
		return cotisation
	case filename == "Sigfaible_delais.csv":
		return delai
	case filename == "Sigfaible_ccsf.csv":
		return ccsf
	case filename == "sireneUL.csv":
		return sireneUl
	case filename == "StockEtablissement_utf8_geo.csv":
		return sirene
	case strings.HasSuffix(filename, ".sas7bdat"):
		return interim
	case mentionsDebits.MatchString(filename):
		return debit
	case hasDianePrefix.MatchString(filename):
		return diane
	case mentionsEffectif.MatchString(filename):
		return effectif
	case hasFilterPrefix.MatchString(filename):
		return filter
	default:
		return ""
	}
}

// These constants represent types supported by our data integration process.
// See https://github.com/signaux-faibles/documentation/blob/master/processus-traitement-donnees.md#sp%C3%A9cificit%C3%A9s-de-limport
const (
	adminUrssaf ValidFileType = "admin_urssaf"
	apconso     ValidFileType = "apconso"
	apdemande   ValidFileType = "apdemande"
	bdf         ValidFileType = "bdf"
	ccsf        ValidFileType = "ccsf"
	cotisation  ValidFileType = "cotisation"
	debit       ValidFileType = "debit"
	delai       ValidFileType = "delai"
	diane       ValidFileType = "diane"
	effectif    ValidFileType = "effectif"
	effectifEnt ValidFileType = "effectif_ent"
	filter      ValidFileType = "filter"
	interim     ValidFileType = "interim"
	procol      ValidFileType = "procol"
	sirene      ValidFileType = "sirene"
	sireneUl    ValidFileType = "sirene_ul"
	paydex      ValidFileType = "paydex"
)

// ValidFileType is the type used by all constants like ADMIN_URSSAF, APCONSO, etc...
type ValidFileType string

var hasDianePrefix = regexp.MustCompile(`^[Dd]iane`)
var mentionsEffectif = regexp.MustCompile(`effectif_`)
var mentionsDebits = regexp.MustCompile(`_debits`)
var hasFilterPrefix = regexp.MustCompile(`^filter_`)
