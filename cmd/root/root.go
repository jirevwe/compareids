package root

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Database connection parameters
	host     string
	port     int
	user     string
	password string
	dbname   string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "compareids",
	Short: "A tool to compare different ID generation strategies",
	Long: `A tool to compare different ID generation strategies for database primary keys.
It can generate test data for each ID type and measure performance metrics.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Define persistent flags for database connection
	RootCmd.PersistentFlags().StringVar(&host, "host", "localhost", "Database host")
	RootCmd.PersistentFlags().IntVar(&port, "port", 5432, "Database port")
	RootCmd.PersistentFlags().StringVar(&user, "user", "postgres", "Database user")
	RootCmd.PersistentFlags().StringVar(&password, "password", "postgres", "Database password")
	RootCmd.PersistentFlags().StringVar(&dbname, "dbname", "postgres", "Database name")
}

// GetDBConnString returns the database connection string
func GetDBConnString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, host, port, dbname)
}
