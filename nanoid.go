package main

import "github.com/matoous/go-nanoid/v2"

// NanoIDGenerator generates NanoIDs
type NanoIDGenerator struct{}

func NewNanoIDGenerator() NanoIDGenerator {
	return NanoIDGenerator{}
}

func (g NanoIDGenerator) Generate() string {
	id, err := gonanoid.New()
	if err != nil {
		panic(err)
	}
	return id
}
