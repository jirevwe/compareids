package common

// TestResult holds the result of a test
type TestResult struct {
	IDType   string
	Count    uint64
	Duration float64
	Stats    map[string]string
}

// TemplateData represents the structure of the data.json file
type TemplateData struct {
	Data      map[string][]map[string]interface{} `json:"Data"`
	IDTypes   []string                            `json:"IDTypes"`
	RowCounts []uint64                            `json:"RowCounts"`
}

// ResultsDir is the directory where individual test results are stored
const ResultsDir = "results"

// DataFile is the path to the template data file
const DataFile = "data.json"
