package core

import (
	"encoding/json"
	"log"
)

// AdminObject represents a document going to be stored in the Admin db collection.
type AdminObject map[string]interface{}

func (this *AdminObject) ToJSON() string {
	json, err := json.MarshalIndent(this, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	return string(json)
}
