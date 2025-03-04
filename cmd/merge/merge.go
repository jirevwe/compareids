package merge

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/jirevwe/compareids/cmd/common"
	"github.com/jirevwe/compareids/cmd/root"
	"github.com/spf13/cobra"
)

// Command represents the merge command
var Command = &cobra.Command{
	Use:   "merge",
	Short: "Merge all test results into a single template_data.json file",
	Long: `Merge all test results from the results directory into a single template_data.json file.
This command should be run after generating test data for all ID types.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load all test results
		results, err := common.LoadTestResults()
		if err != nil {
			log.Fatalf("Error loading test results: %v", err)
		}

		if len(results) == 0 {
			log.Fatalf("No test results found in %s directory", common.ResultsDir)
		}

		// Create a map to store the row counts
		rowCountsMap := make(map[uint64]bool)
		for _, result := range results {
			rowCountsMap[result.Count] = true
		}

		// Convert the map to a sorted slice
		rowCounts := make([]uint64, 0, len(rowCountsMap))
		for count := range rowCountsMap {
			rowCounts = append(rowCounts, count)
		}
		sort.Slice(rowCounts, func(i, j int) bool {
			return rowCounts[i] < rowCounts[j]
		})

		// Create a map to store the ID types
		idTypesMap := make(map[string]bool)
		for _, result := range results {
			idTypesMap[result.IDType] = true
		}

		// Convert the map to a sorted slice
		idTypes := make([]string, 0, len(idTypesMap))
		for idType := range idTypesMap {
			idTypes = append(idTypes, idType)
		}
		sort.Strings(idTypes)

		// Create the template data
		templateData := common.TemplateData{
			Data:      make(map[string][]map[string]interface{}),
			IDTypes:   idTypes,
			RowCounts: rowCounts,
		}

		// Process the results for the template
		for _, idType := range idTypes {
			templateData.Data[idType] = make([]map[string]interface{}, 0)

			// Get all results for this ID type
			for _, result := range results {
				if result.IDType == idType {
					// Create a data point with the count and duration
					dataPoint := map[string]interface{}{
						"count":    result.Count,
						"duration": result.Duration,
					}

					// Add the stats directly to the data point
					// This ensures that the HTML file can access them directly
					for k, v := range result.Stats {
						// Don't overwrite the count field
						if k != "count" {
							dataPoint[k] = v
						}
					}

					// Make sure all required fields exist
					requiredFields := []string{
						"total_table_size", "data_size", "index_size",
						"index_internal_pages", "index_leaf_pages",
						"index_density", "index_fragmentation",
						"index_internal_to_leaf_ratio",
					}

					for _, field := range requiredFields {
						if _, exists := dataPoint[field]; !exists {
							// If a field is missing, add a default value
							dataPoint[field] = "0"
						}
					}

					templateData.Data[idType] = append(templateData.Data[idType], dataPoint)
				}
			}

			// Sort the data points by count
			sort.Slice(templateData.Data[idType], func(i, j int) bool {
				// Safely handle the count value which might be a string or uint64
				var countI, countJ uint64

				switch c := templateData.Data[idType][i]["count"].(type) {
				case uint64:
					countI = c
				case string:
					// Convert string to uint64
					countI, _ = strconv.ParseUint(c, 10, 64)
				}

				switch c := templateData.Data[idType][j]["count"].(type) {
				case uint64:
					countJ = c
				case string:
					// Convert string to uint64
					countJ, _ = strconv.ParseUint(c, 10, 64)
				}

				return countI < countJ
			})
		}

		// Write template data to JSON file
		templateDataFile, err := os.Create(common.TemplateDataFile)
		if err != nil {
			log.Fatalf("Unable to create template data JSON file: %v\n", err)
		}
		defer templateDataFile.Close()

		encoder := json.NewEncoder(templateDataFile)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(templateData); err != nil {
			log.Fatalf("Unable to write template data to JSON file: %v\n", err)
		}

		fmt.Printf("Merged %d test results into %s\n", len(results), common.TemplateDataFile)
	},
}

func init() {
	// Add the merge command to the root command
	root.RootCmd.AddCommand(Command)
}
