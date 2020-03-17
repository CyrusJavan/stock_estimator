package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "stock_estimator",
	Short: "Simulate investment strategy performance",
	Long: `Simulate investment strategy performance using historical data.
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf(`Usage: stock_estimator invest [flags]
`)
	},
}

func init() {
	investCmd.PersistentFlags().Float64Var(&Principal, "principal", 0, "inital principal")
	investCmd.PersistentFlags().Float64Var(&RecurringInvestment, "recur", 0, "recurring investment")
	investCmd.PersistentFlags().StringVar(&StartDate, "start", "1985-01-01", "start date")
	investCmd.PersistentFlags().StringVar(&EndDate, "end", "2020-01-01", "end date")

	rootCmd.AddCommand(investCmd)
	rootCmd.AddCommand(serverCmd)
}

// Execute is the entry point for the cobra app
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
