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

// PopulateAdminObject populates an AdminObject, given a list of data files.
func PopulateAdminObject(augmentedFilenames []DataFile, batchKey BatchKey, dateFinEffectif DateFinEffectif) (AdminObject, UnsupportedFilesError) {

	filesProperty, unsupportedFiles := PopulateFilesProperty(augmentedFilenames, batchKey.Path())
	var completeTypes = []ValidFileType{}
	for _, typeName := range defaultCompleteTypes {
		if _, ok := filesProperty[typeName]; ok {
			completeTypes = append(completeTypes, typeName)
		}
	}
	// { "date_debut" : { "$date" : "2014-01-01T00:00:00.000+0000" }, "date_fin" : { "$date" : "2018-12-01T00:00:00.000+0000" }, "date_fin_effectif" : { "$date" : "2018-06-01T00:00:00.000+0000" } }

	paramProperty := ParamProperty{
		DateDebut:       MongoDate{"2014-01-01T00:00:00.000+0000"},
		DateFin:         MongoDate{"20" + batchKey.String()[0:2] + "-" + batchKey.String()[2:4] + "-01T00:00:00.000+0000"},
		DateFinEffectif: dateFinEffectif.MongoDate(),
	}

	return AdminObject{
		"_id":            IDProperty{batchKey, "batch"},
		"files":          filesProperty,
		"complete_types": completeTypes,
		"param":          paramProperty,
	}, UnsupportedFilesError{unsupportedFiles}
}

// UnsupportedFilesError is an Error object that lists files that were not supported.
type UnsupportedFilesError struct {
	UnsupportedFiles []string
}

func (err UnsupportedFilesError) Error() string {
	return "unsupported: " + strings.Join(err.UnsupportedFiles, ", ")
}

var defaultCompleteTypes = []ValidFileType{
	apconso,
	apdemande,
	effectif,
	effectifEnt,
	sirene,
	sireneUl,
}
