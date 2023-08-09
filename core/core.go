// Package core contient le code commun à tous les packages. Il doit contenir le code métier de l'application
package core

import (
	"encoding/json"
	"log"
)

// AdminObject represents a document going to be stored in the Admin db collection.
type AdminObject map[string]interface{}

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
