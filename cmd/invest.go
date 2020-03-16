package cmd

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"time"

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

const (
	timeFormat string = "2006-01-02"
)

func invest(cmd *cobra.Command, args []string) {
	f, err := os.Open("data/DJI.csv")
	if err != nil {
		panic(err)
	}
	r := csv.NewReader(f)
	dji := stock{}
	r.Read() // Read the header line

	for {
		line, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		// Reference time: Mon Jan 2 15:04:05 -0700 MST 2006
		val, err := strconv.ParseFloat(line[1], 32)
		if err != nil {
			panic(err)
		}
		d, err := time.Parse(timeFormat, line[0])
		if err != nil {
			panic(err)
		}
		dji.prices = append(dji.prices, price{val, d})
	}

	start, _ := time.Parse(timeFormat, StartDate)
	end, _ := time.Parse(timeFormat, EndDate)

	p := message.NewPrinter(language.English)

	worth, invested := dji.investOverTime(
		start,
		end,
		Principal,
		RecurringInvestment,
	)

	p.Printf(
		`Initial Principal:    $%10.2f
Recurring Investment: $%10.2f
Start Date:           %s
End Date:             %s
Total Invested:       $%10.2f
Resulting Net Worth:  $%10.2f
Change:               %10.2f%%
`, Principal, RecurringInvestment, StartDate, EndDate, invested, worth, ((worth/invested)-1)*100)
}

type price struct {
	value float64
	date  time.Time
}

type stock struct {
	prices []price
}

func (s stock) investOverTime(
	start,
	end time.Time,
	principal,
	recurringInvestment float64) (float64, float64) {
	var beginPrice, endPrice, shares, recurringTotal float64
	for _, pricePoint := range s.prices {
		if pricePoint.date.Equal(start) || (pricePoint.date.After(start) && beginPrice == 0) {
			beginPrice = pricePoint.value
			shares = principal / beginPrice
		} else if pricePoint.date.After(start) && pricePoint.date.Before(end) {
			shares += recurringInvestment / pricePoint.value
			recurringTotal += recurringInvestment
		} else if pricePoint.date.Equal(end) || (pricePoint.date.After(end) && beginPrice != 0) {
			endPrice = pricePoint.value
			break
		}
	}
	return shares * endPrice, principal + recurringTotal
}
