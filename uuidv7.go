package main

import (
	"github.com/gofrs/uuid/v5"
)

// UUIDv7Generator generates UUIDv7 IDs
type UUIDv7Generator struct{}

func NewUUIDv7Generator() UUIDv7Generator {
	return UUIDv7Generator{}
}

func (g UUIDv7Generator) Generate() string {
	id, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	return id.String()
}
