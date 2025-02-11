package main

import (
	"github.com/oklog/ulid/v2"
)

// ULIDGenerator generates ULID IDs
type ULIDGenerator struct{}

func NewULIDGenerator() ULIDGenerator {
	return ULIDGenerator{}
}

func (g ULIDGenerator) Generate() string {
	return ulid.Make().String()
}
