package main

import (
	"regexp"
	"strings"
)

var hasDianePrefix = regexp.MustCompile(`^diane`)
var mentionsEffectif = regexp.MustCompile(`effectif_`)
var mentionsDebits = regexp.MustCompile(`_debits`)
var hasFilterPrefix = regexp.MustCompile(`^filter_`)

type MetadataProperty map[string]string

type UploadedFileMeta struct {
	MetaData MetadataProperty
}

func ExtractFileTypeFromMetadata(filename string, fileinfo UploadedFileMeta) string {
	metadata := fileinfo.MetaData
	if metadata["goup-path"] == "bdf" {
		return "bdf"
	} else {
		return ExtractFileTypeFromFilename(metadata["filename"])
	}
}

// ExtractFileTypeFromFilename returns a file type from filename, or empty string for unsupported file names
func ExtractFileTypeFromFilename(filename string) string {
	switch {
	case filename == "act_partielle_conso_depuis2014_FRANCE.csv":
		return "apconso"
	case filename == "act_partielle_ddes_depuis2015_FRANCE.csv":
		return "apdemande"
	case filename == "Sigfaible_etablissement_utf8.csv":
		return "admin_urssaf"
	case filename == "Sigfaible_effectif_siren.csv":
		return "effectif_ent"
	case filename == "Sigfaible_pcoll.csv":
		return "procol"
	case filename == "Sigfaible_cotisdues.csv":
		return "cotisation"
	case filename == "Sigfaible_delais.csv":
		return "delai"
	case filename == "Sigfaible_ccsf.csv":
		return "ccsf"
	case filename == "sireneUL.csv":
		return "sirene_ul"
	case filename == "StockEtablissement_utf8_geo.csv":
		return "comptes"
	case strings.HasSuffix(filename, ".sas7bdat"):
		return "interim"
	case mentionsDebits.MatchString(filename):
		return "debit"
	case hasDianePrefix.MatchString(filename):
		return "diane"
	case mentionsEffectif.MatchString(filename):
		return "effectif"
	case hasFilterPrefix.MatchString(filename):
		return "filter"
	default:
		return ""
	}
}
