package cmd

import (
	"github.com/CyrusJavan/stock_estimator/simulation"
	"github.com/spf13/cobra"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var investCmd = &cobra.Command{
	Use:   "invest",
	Short: "Simulate investment strategy performance",
	Long:  `Simulate investment strategy performance using historical data.`,
	Run:   invest,
}

var (
	// Principal is how much money we start with
	Principal float64

	// RecurringInvestment is how much we will invest every trading day
	RecurringInvestment float64

	// StartDate is when we will begin investing
	StartDate string

	// EndDate is when we end investing
	EndDate string
)

func invest(cmd *cobra.Command, args []string) {
	sim, err := simulation.NewSimulation("data/DJI.csv")

	if err != nil {
		panic(err)
	}

	worth, invested := sim.InvestOverTime(
		StartDate,
		EndDate,
		Principal,
		RecurringInvestment,
	)

	p := message.NewPrinter(language.English)

	p.Printf(
		`Initial Principal:    $%10.2f
Recurring Investment: $%10.2f
Start Date:           %s
End Date:             %s
Total Invested:       $%10.2f
Resulting Net Worth:  $%10.2f
Change:               %10.2f%%
`, Principal, RecurringInvestment, StartDate, EndDate,
		invested, worth, ((worth/invested)-1)*100)
}
