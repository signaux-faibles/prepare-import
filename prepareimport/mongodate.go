package prepareimport

import "encoding/json"

// MongoDate represents a date that can be serialized for MongoDB.
type MongoDate struct {
	date string
}

// MarshalJSON will be called when serializing a date for MongoDB.
func (mongoDate MongoDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{"$date": mongoDate.date})
}

// UnmarshalJSON will parse a MongoDate from a MongoDB date object.
func (mongoDate *MongoDate) UnmarshalJSON(data []byte) error {
	var dateObj map[string]string
	if err := json.Unmarshal(data, &dateObj); err != nil {
		return err
	}
	mongoDate.date = dateObj["$date"]
	return nil
}
