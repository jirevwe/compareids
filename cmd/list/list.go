package list

import (
	"fmt"

	"github.com/jirevwe/compareids/cmd/common"
	"github.com/jirevwe/compareids/cmd/root"
	"github.com/spf13/cobra"
)

// Command represents the list command
var Command = &cobra.Command{
	Use:   "list",
	Short: "List all available ID types",
	Long:  `List all available ID types that can be used with the id command.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get all ID types
		idTypes := common.GetAllIDTypes()

		// Print the ID types
		fmt.Println("Available ID types:")
		for _, idType := range idTypes {
			fmt.Printf("  - %s\n", idType)
		}
	},
}

func init() {
	// Add the list command to the root command
	root.RootCmd.AddCommand(Command)
}
