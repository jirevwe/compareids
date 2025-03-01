package main

import (
	// Import the commands
	_ "github.com/jirevwe/compareids/cmd/all"
	_ "github.com/jirevwe/compareids/cmd/id"
	_ "github.com/jirevwe/compareids/cmd/list"
	_ "github.com/jirevwe/compareids/cmd/merge"
	"github.com/jirevwe/compareids/cmd/root"
)

func main() {
	// Execute the root command
	root.Execute()
}
