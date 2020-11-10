package prepareimport

import (
	"strings"
)

// AdminObject represents a document going to be stored in the Admin db collection.
type AdminObject map[string]interface{}

// IDProperty represents the "_id" property of an Admin object.
type IDProperty struct {
	Key  BatchKey `json:"key"`
	Type string   `json:"type"`
}

// UnsupportedFilesError is an Error object that lists files that were not supported.
type UnsupportedFilesError struct {
	UnsupportedFiles []string
}

func (err UnsupportedFilesError) Error() string {
	return "unsupported: " + strings.Join(err.UnsupportedFiles, ", ")
}

func populateParamProperty(batchKey BatchKey, dateFinEffectif DateFinEffectif) ParamProperty {
	return ParamProperty{
		DateDebut:       MongoDate{"2014-01-01T00:00:00.000+0000"},
		DateFin:         MongoDate{"20" + batchKey.String()[0:2] + "-" + batchKey.String()[2:4] + "-01T00:00:00.000+0000"},
		DateFinEffectif: dateFinEffectif.MongoDate(),
	}
}

func populateCompleteTypesProperty(filesProperty FilesProperty) []ValidFileType {
	completeTypes := []ValidFileType{}
	for _, typeName := range defaultCompleteTypes {
		if _, ok := filesProperty[typeName]; ok {
			completeTypes = append(completeTypes, typeName)
		}
	}
	return completeTypes
}

var defaultCompleteTypes = []ValidFileType{
	apconso,
	apdemande,
	effectif,
	effectifEnt,
	sirene,
	sireneUl,
}
