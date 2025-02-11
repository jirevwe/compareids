package main

import (
	"github.com/bwmarrin/snowflake"
	"log"
)

// SnowflakeGenerator generates Snowflake IDs
type SnowflakeGenerator struct {
	node *snowflake.Node
}

func NewSnowflakeGenerator() SnowflakeGenerator {
	node, err := snowflake.NewNode(1)
	if err != nil {
		log.Fatalf("Failed to create Snowflake node: %v", err)
	}
	return SnowflakeGenerator{node: node}
}

func (g SnowflakeGenerator) Generate() string {
	return g.node.Generate().String()
}
