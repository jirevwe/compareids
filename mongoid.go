package main

import "go.mongodb.org/mongo-driver/bson/primitive"

// MongoIDGenerator generates MongoDB ObjectIDs
type MongoIDGenerator struct{}

func NewMongoIDGenerator() MongoIDGenerator {
	return MongoIDGenerator{}
}

func (g MongoIDGenerator) Generate() string {
	return primitive.NewObjectID().Hex()
}
