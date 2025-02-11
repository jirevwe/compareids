package main

import "github.com/rs/xid"

// XIDGenerator generates XIDs
type XIDGenerator struct{}

func NewXIDGenerator() XIDGenerator {
	return XIDGenerator{}
}

func (g XIDGenerator) Generate() string {
	return xid.New().String()
}
