// Package core contient le code commun à tous les packages. Il doit contenir le code métier de l'application
package core

import (
	"bytes"
	"encoding/json"
	"log"
	"time"
)

// AdminObject represents a document going to be stored in the Admin db collection.
//type AdminObject map[string]interface{}

type AdminObject struct {
	ID            interface{}          `bson:"_id"`
	CompleteTypes []string             `bson:"complete_types"`
	Files         map[string][]string  `bson:"files"`
	Params        map[string]time.Time `bson:"param"`
}

// ToJSON retourne le json caractérisant l'objet passé en paramètre
func (current *AdminObject) ToJSON() string {
	jsonText, err := json.MarshalIndent(current, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	return string(jsonText)
}

// ToCompactJSON retourne le json caractérisant l'objet passé en paramètre
func (current *AdminObject) ToCompactJSON() string {
	compact := &bytes.Buffer{}
	if err := json.Compact(compact, []byte(current.ToJSON())); err != nil {
		panic(err)
	}
	return compact.String()
}

func FromJSON(input string) AdminObject {
	var r AdminObject
	err := json.Unmarshal([]byte(input), &r)
	if err != nil {
		panic(err)
	}
	return r
}
