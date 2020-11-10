package prepareimport

import (
	"strings"
	"time"
)

// DateFinEffectif is a date that can be serialized for MongoDB.
type DateFinEffectif interface {
	MongoDate() MongoDate
}

// NewDateFinEffectif creates a valid DateFinEffectif.
func NewDateFinEffectif(date time.Time) DateFinEffectif {
	// TODO: validate the value
	return dateFinEffectifType{date}
}

type dateFinEffectifType struct {
	time.Time
}

func (dateFinEffectif dateFinEffectifType) MongoDate() MongoDate {
	return MongoDate{strings.Replace(dateFinEffectif.Format(time.RFC3339), "Z", ".000+0000", 1)}
}
