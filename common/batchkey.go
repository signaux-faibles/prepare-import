package common

import (
	"errors"
	"regexp"
)

type batchKeyType string

func (b batchKeyType) String() string {
	return string(b)
}

func (b batchKeyType) Path() string {
	return "/" + string(b) + "/"
}

// BatchKey represents a valid batch key.
type BatchKey interface {
	String() string
	Path() string
}

// NewBatchKey constructs a valid batch key.
func NewBatchKey(key string) (BatchKey, error) {
	var isValidBatchKey = regexp.MustCompile(`^[0-9]{4}`)
	if !isValidBatchKey.MatchString(key) {
		return batchKeyType(""), errors.New("la clé du batch doit respecter le format requis AAMM")
	}
	return batchKeyType(key), nil
}
