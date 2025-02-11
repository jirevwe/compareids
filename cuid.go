package main

import "github.com/lucsky/cuid"

// CUIDGenerator generates CUIDs
type CUIDGenerator struct{}

func NewCUIDGenerator() CUIDGenerator {
	return CUIDGenerator{}
}

func (g CUIDGenerator) Generate() string {
	return cuid.New()
}
