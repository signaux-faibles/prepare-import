package prepareimport

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

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

// ToTime
func (mongoDate MongoDate) ToTime() (time.Time, error) {
	layout := "2006-01-02T15:04:05.000+0000"
	date, err := time.Parse(layout, mongoDate.date)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "erreur pendant le parsing de la date mongo : "+mongoDate.date)
	}
	return date, nil
}
