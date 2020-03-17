package simulation

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"time"
)

const (
	timeFormat string = "2006-01-02"
)

// Simulation is the type used to perform investment sims
type Simulation struct {
	s stock
}

// NewSimulation makes a new simulation based on the data file
func NewSimulation(dataFile string) Simulation {
	f, err := os.Open(dataFile)
	if err != nil {
		panic(err)
	}
	r := csv.NewReader(f)
	s := stock{}
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
		s.prices = append(s.prices, price{val, d})
	}
	return Simulation{s}
}

// InvestOverTime returns the amount of money earned and amount invested
func (s Simulation) InvestOverTime(
	startDate,
	endDate string,
	principal,
	recurringInvestment float64) (float64, float64) {
	start, _ := time.Parse(timeFormat, startDate)
	end, _ := time.Parse(timeFormat, endDate)

	return s.s.investOverTime(start, end, principal, recurringInvestment)
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
