package prepareimport

import (
	"errors"
	"regexp"
)

// BatchKey represents a valid batch key.
type BatchKey interface {
	String() string
	Path() string
	IsSubBatch() bool
	GetParentBatch() string
}

// NewBatchKey constructs a valid batch key.
func NewBatchKey(key string) (BatchKey, error) {
	if !validBatchKey.MatchString(key) {
		return batchKeyType(""), errors.New("la clé du batch doit respecter le format requis AAMM")
	}
	return batchKeyType(key), nil
}

var validBatchKey = regexp.MustCompile(`^[0-9]{4}`)
var validSubBatchKey = regexp.MustCompile(`^([0-9]{4})_([0-9]{2})$`)

type batchKeyType string

func (b batchKeyType) String() string {
	return string(b)
}

func (b batchKeyType) Path() string {
	return "/" + string(b) + "/"
}

func (b batchKeyType) IsSubBatch() bool {
	return validSubBatchKey.MatchString(string(b))
}

func (b batchKeyType) GetParentBatch() string {
	matches := validSubBatchKey.FindStringSubmatch(string(b))
	return matches[1]
}
