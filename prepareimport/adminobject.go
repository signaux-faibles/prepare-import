package prepareimport

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"prepare-import/core"
)

// IDProperty represents the "_id" property of an Admin object.
type IDProperty struct {
	Key  BatchKey `bson:"key,omitempty" json:"key,omitempty"`
	Type string   `bson:"type,omitempty" json:"type,omitempty"`
}

// ParamProperty represents the "param" property of an Admin object.
type ParamProperty struct {
	DateDebut       time.Time `bson:"date_debut" json:"date_debut"`
	DateFin         time.Time `bson:"date_fin" json:"date_fin"`
	DateFinEffectif time.Time `bson:"date_fin_effectif" json:"date_fin_effectif"`
}

// UnsupportedFilesError is an Error object that lists files that were not supported.
type UnsupportedFilesError struct {
	UnsupportedFiles []string
}

func (err UnsupportedFilesError) Error() string {
	return "type de fichier non supporté : " + strings.Join(err.UnsupportedFiles, ", ")
}

func populateParamProperty(batchKey BatchKey, dateFinEffectif DateFinEffectif) ParamProperty {
	year, _ := strconv.Atoi("20" + batchKey.String()[0:2])
	month, _ := strconv.Atoi(batchKey.String()[2:4])
	return ParamProperty{
		DateDebut:       time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC),
		DateFin:         time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC),
		DateFinEffectif: dateFinEffectif.Date(),
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

func populateFilesPaths(filesProperty FilesProperty) map[ValidFileType][]string {
	r := make(map[ValidFileType][]string)
	for k, v := range filesProperty {
		r[k] = core.Apply(v, func(bf BatchFile) string { return bf.Path() })
	}
	return r
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
