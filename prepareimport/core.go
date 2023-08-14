package prepareimport

import (
	"encoding/json"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

// AdminObject represents a document going to be stored in the Admin db collection.
type AdminObject struct {
	ID            IDProperty                 `bson:"_id,omitempty" json:"_id,omitempty"`
	CompleteTypes []ValidFileType            `bson:"complete_types,omitempty" json:"complete_types,omitempty"`
	Files         map[ValidFileType][]string `bson:"files,omitempty" json:"files,omitempty"`
	Param         ParamProperty              `bson:"param,omitempty" json:"param,omitempty"`
}

// ToJSON retourne le json caractérisant l'objet passé en paramètre
func (current *AdminObject) ToJSON() string {
	jsonText, err := json.MarshalIndent(current, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	return string(jsonText)
}

func FromJSON(input string) AdminObject {
	var r AdminObject
	err := json.Unmarshal([]byte(input), &r)
	if err != nil {
		panic(err)
	}
	return r
}

func FromBSON(input []byte) AdminObject {
	var r AdminObject
	//var temp interface{}
	err := bson.Unmarshal(input, &r)
	if err != nil {
		panic(err)
	}
	return r
}
