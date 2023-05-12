package prepareimport

import (
	"fmt"
	"strings"
)

// IDProperty represents the "_id" property of an Admin object.
type IDProperty struct {
	Key  BatchKey `json:"key"`
	Type string   `json:"type"`
}

// ParamProperty represents the "param" property of an Admin object.
type ParamProperty struct {
	DateDebut       MongoDate `json:"date_debut"`
	DateFin         MongoDate `json:"date_fin"`
	DateFinEffectif MongoDate `json:"date_fin_effectif"`
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
		DateDebut:       MongoDate{"2016-01-01T00:00:00.000+0000"},
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
	for typeName, thresholdInBytes := range thresholdPerGzippedFileType {
		if files, ok := filesProperty[typeName]; ok {
			if len(files) != 1 {
				panic(fmt.Errorf("'complete' file detection can only work if there is only 1 file per type, found %v for type %v", len(files), typeName))
			}
			file := files[0]
			if file.GetGzippedSize() >= thresholdInBytes {
				println(fmt.Sprintf("Info: file \"%v\" was marked as \"complete\" because it's a gzipped file which size reached the threshold of %v bytes", file.Name(), thresholdInBytes))
				completeTypes = append(completeTypes, typeName)
			}
		}
	}
	return completeTypes
}

// types of files that are always provided as "complete"
var defaultCompleteTypes = []ValidFileType{
	apconso,
	apdemande,
	effectif,
	effectifEnt,
	sirene,
	sireneUl,
}

// types of files that will be considered as "complete" if their gzipped size reach a certain threshold (in bytes)
var thresholdPerGzippedFileType = map[ValidFileType]uint64{
	cotisation: 143813078,
	delai:      1666199,
	procol:     1646193,
	debit:      254781489,
}
