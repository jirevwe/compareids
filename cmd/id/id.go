package id

import (
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jirevwe/compareids/cmd/common"
	"github.com/jirevwe/compareids/cmd/root"
	"github.com/spf13/cobra"
)

var (
	// rowCount is the number of rows to generate
	rowCount uint64
)

// Command represents the id command
var Command = &cobra.Command{
	Use:   "id [id-type]",
	Short: "Generate test data for a specific ID type",
	Long: `Generate test data for a specific ID type and save the results to a JSON file.
Example: compareids id uuidv4 --count 10000`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		idType := args[0]

		// Get the ID generator
		generator, err := common.GetIDGenerator(idType)
		if err != nil {
			log.Fatalf("Error getting ID generator: %v", err)
		}

		// Create a connection pool
		connString := root.GetDBConnString()
		config, err := pgxpool.ParseConfig(connString)
		if err != nil {
			log.Fatalf("Unable to parse connection string: %v\n", err)
		}

		pool, err := pgxpool.NewWithConfig(cmd.Context(), config)
		if err != nil {
			log.Fatalf("Unable to create connection pool: %v\n", err)
		}
		defer pool.Close()

		// Run the test
		fmt.Printf("Running test for %s with %d rows...\n", generator.Name(), rowCount)
		duration, stats, err := common.RunTest(cmd.Context(), pool, generator, rowCount)
		if err != nil {
			log.Fatalf("Error running test: %v", err)
		}

		// Drop the table
		if err = generator.DropTable(cmd.Context(), pool); err != nil {
			log.Printf("Error dropping table: %v", err)
		}

		// Save the result
		result := common.TestResult{
			IDType:   generator.Name(),
			Count:    rowCount,
			Duration: duration,
			Stats:    stats,
		}

		if err := common.SaveTestResult(result); err != nil {
			log.Fatalf("Error saving test result: %v", err)
		}

		fmt.Printf("Test completed in %.2fms. Results saved to %s/%s_%d.json\n",
			duration, common.ResultsDir, generator.Name(), rowCount)
	},
}

func init() {
	// Add the id command to the root command
	root.RootCmd.AddCommand(Command)

	// Define flags
	Command.Flags().Uint64Var(&rowCount, "count", 10000, "Number of rows to generate")
}

// GetSupportedIDTypes returns a list of supported ID types
func GetSupportedIDTypes() []string {
	return common.GetAllIDTypes()
}
