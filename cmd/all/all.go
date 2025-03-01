package all

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jirevwe/compareids/cmd/common"
	"github.com/jirevwe/compareids/cmd/merge"
	"github.com/jirevwe/compareids/cmd/root"
	"github.com/spf13/cobra"
)

var (
	// skipMerge is a flag to skip merging the results
	skipMerge bool
)

// Command represents the all command
var Command = &cobra.Command{
	Use:   "all",
	Short: "Run tests for all ID types",
	Long: `Run tests for all ID types with the default row counts and merge the results.
This is equivalent to running the id command for each ID type and then the merge command.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get all ID types
		idTypes := common.GetAllIDTypes()
		rowCounts := common.GetDefaultRowCounts()

		// Create a connection pool
		connString := root.GetDBConnString()
		config, err := pgxpool.ParseConfig(connString)
		if err != nil {
			log.Fatalf("Unable to parse connection string: %v\n", err)
		}

		pool, err := pgxpool.NewWithConfig(context.Background(), config)
		if err != nil {
			log.Fatalf("Unable to create connection pool: %v\n", err)
		}
		defer pool.Close()

		// Run tests for all ID types and row counts
		for _, idType := range idTypes {
			// Get the ID generator
			generator, err := common.GetIDGenerator(idType)
			if err != nil {
				log.Printf("Error getting ID generator for %s: %v", idType, err)
				continue
			}

			for _, count := range rowCounts {
				// Run the test
				fmt.Printf("Running test for %s with %d rows...\n", generator.Name(), count)
				duration, stats, err := common.RunTest(pool, generator, count)
				if err != nil {
					log.Printf("Error running test for %s with %d rows: %v", generator.Name(), count, err)
					continue
				}

				// Save the result
				result := common.TestResult{
					IDType:   generator.Name(),
					Count:    count,
					Duration: duration,
					Stats:    stats,
				}

				if err := common.SaveTestResult(result); err != nil {
					log.Printf("Error saving test result for %s with %d rows: %v", generator.Name(), count, err)
					continue
				}

				fmt.Printf("Test completed in %.2fms. Results saved to %s/%s_%d.json\n",
					duration, common.ResultsDir, generator.Name(), count)
			}

			// Drop the table
			if err = generator.DropTable(context.Background(), pool); err != nil {
				log.Printf("Error dropping table for %s: %v", generator.Name(), err)
			}
		}

		// Merge the results
		if !skipMerge {
			fmt.Println("Merging results...")
			merge.Command.Run(cmd, args)
		}
	},
}

func init() {
	// Add the all command to the root command
	root.RootCmd.AddCommand(Command)

	// Define flags
	Command.Flags().BoolVar(&skipMerge, "skip-merge", false, "Skip merging the results")
}
