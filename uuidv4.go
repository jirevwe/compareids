package main

import "github.com/google/uuid"

// UUIDv4Generator generates UUIDv4 IDs
type UUIDv4Generator struct{}

func NewUUIDv4Generator() UUIDv4Generator {
	return UUIDv4Generator{}
}

func (g UUIDv4Generator) Generate() string {
	return uuid.New().String()
}
