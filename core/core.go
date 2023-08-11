// Package core contient le code commun à tous les packages. Il doit contenir le code métier de l'application
package core

import (
	"encoding/json"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

// AdminObject represents a document going to be stored in the Admin db collection.
type AdminObject bson.M

// ToJSON retourne le json caractérisant l'objet passé en paramètre
func (current *AdminObject) ToJSON() string {
	jsonText, err := json.MarshalIndent(current, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	return string(jsonText)
}

func (current AdminObject) GetKey() string {
	id := current["_id"].(map[string]interface{})
	key := fmt.Sprint(id["key"])
	return key
}

func FromJSON(input string) AdminObject {
	var r AdminObject
	err := json.Unmarshal([]byte(input), &r)
	if err != nil {
		panic(err)
	}
	return r
}
