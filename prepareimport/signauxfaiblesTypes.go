package prepareimport

import (
	"regexp"
	"strings"
)

// These constants represent types supported by our data integration process.
// See https://github.com/signaux-faibles/documentation/blob/master/processus-traitement-donnees.md#sp%C3%A9cificit%C3%A9s-de-limport
const (
	ADMIN_URSSAF ValidFileType = "admin_urssaf"
	APCONSO      ValidFileType = "apconso"
	APDEMANDE    ValidFileType = "apdemande"
	BDF          ValidFileType = "bdf"
	CCSF         ValidFileType = "ccsf"
	COTISATION   ValidFileType = "cotisation"
	DEBIT        ValidFileType = "debit"
	DELAI        ValidFileType = "delai"
	DIANE        ValidFileType = "diane"
	EFFECTIF     ValidFileType = "effectif"
	EFFECTIF_ENT ValidFileType = "effectif_ent"
	FILTER       ValidFileType = "filter"
	INTERIM      ValidFileType = "interim"
	PROCOL       ValidFileType = "procol"
	SIRENE       ValidFileType = "sirene"
	SIRENE_UL    ValidFileType = "sirene_ul"
)

// ValidFileType is the type used by all constants like ADMIN_URSSAF, APCONSO, etc...
type ValidFileType string

var defaultCompleteTypes = []ValidFileType{
	APCONSO,
	APDEMANDE,
	EFFECTIF,
	EFFECTIF_ENT,
	SIRENE,
	SIRENE_UL,
}

var hasDianePrefix = regexp.MustCompile(`^[Dd]iane`)
var mentionsEffectif = regexp.MustCompile(`effectif_`)
var mentionsDebits = regexp.MustCompile(`_debits`)
var hasFilterPrefix = regexp.MustCompile(`^filter_`)

// MetadataProperty represents the "MetaData" property from .info files.
type MetadataProperty map[string]string

// UploadedFileMeta represents the JSON object loaded from .info files.
type UploadedFileMeta struct {
	MetaData MetadataProperty
}

// ExtractFileTypeFromMetadata returns the type of a bin file (without extension), based on the contents of the associated .info file.
func ExtractFileTypeFromMetadata(filename string, fileinfo UploadedFileMeta) ValidFileType {
	metadata := fileinfo.MetaData
	if metadata["goup-path"] == "bdf" {
		return BDF
	}
	return ExtractFileTypeFromFilename(metadata["filename"])
}

// ExtractFileTypeFromFilename returns a file type from filename, or empty string for unsupported file names
func ExtractFileTypeFromFilename(filename string) ValidFileType {
	switch {
	case filename == "act_partielle_conso_depuis2014_FRANCE.csv":
		return APCONSO
	case filename == "act_partielle_ddes_depuis2015_FRANCE.csv":
		return APDEMANDE
	case filename == "Sigfaible_etablissement_utf8.csv":
		return ADMIN_URSSAF
	case filename == "Sigfaible_effectif_siren.csv":
		return EFFECTIF_ENT
	case filename == "Sigfaible_pcoll.csv":
		return PROCOL
	case filename == "Sigfaible_cotisdues.csv":
		return COTISATION
	case filename == "Sigfaible_delais.csv":
		return DELAI
	case filename == "Sigfaible_ccsf.csv":
		return CCSF
	case filename == "sireneUL.csv":
		return SIRENE_UL
	case filename == "StockEtablissement_utf8_geo.csv":
		return SIRENE
	case strings.HasSuffix(filename, ".sas7bdat"):
		return INTERIM
	case mentionsDebits.MatchString(filename):
		return DEBIT
	case hasDianePrefix.MatchString(filename):
		return DIANE
	case mentionsEffectif.MatchString(filename):
		return EFFECTIF
	case hasFilterPrefix.MatchString(filename):
		return FILTER
	default:
		return ""
	}
}
