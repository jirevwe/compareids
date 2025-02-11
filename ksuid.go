package main

import "github.com/segmentio/ksuid"

// KSUIDGenerator generates KSUIDs
type KSUIDGenerator struct{}

func NewKSUIDGenerator() KSUIDGenerator {
	return KSUIDGenerator{}
}

func (g KSUIDGenerator) Generate() string {
	return ksuid.New().String()
}
