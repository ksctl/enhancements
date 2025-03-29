package main

import (
	"fmt"
	"log"

	"github.com/ksctl/enhancements/poc/co2_intensity/climatetrace"
	"github.com/ksctl/enhancements/poc/co2_intensity/electricitymaps"
	"github.com/ksctl/enhancements/poc/co2_intensity/ember"
)

func handleClimateTrace() {
	ct := climatetrace.DefaultClimateTrace

	if v, err := ct.GetCountries(); err != nil {
		log.Fatal(err)
	} else {
		v.P()
	}

	if v, err := ct.GetEmissionSummaryHistroy(); err != nil {
		log.Fatal(err)
	} else {
		v.P()
	}
}

func handleElectricityMaps() {
	em := electricitymaps.DefaultElectricityMap

	if v, err := em.GetAvailableZones(); err != nil {
		log.Fatal(err)
	} else {
		v.S()
	}

	if v, err := em.GetMonthlyPastData("IN-EA"); err != nil {
		log.Fatal(err)
	} else {
		v.S()
	}
	if v, err := em.GetMonthlyPastData("IN-SO"); err != nil {
		log.Fatal(err)
	} else {
		v.S()
	}
	if v, err := em.GetMonthlyPastData("IN-WE"); err != nil {
		log.Fatal(err)
	} else {
		v.S()
	}
	if v, err := em.GetMonthlyPastData("IN-NE"); err != nil {
		log.Fatal(err)
	} else {
		v.S()
	}
	if v, err := em.GetMonthlyPastData("IN-NO"); err != nil {
		log.Fatal(err)
	} else {
		v.S()
	}

	if v, err := em.GetLatestCarbonIntensity("IN-EA", electricitymaps.OptionEmissionFactorType("direct")); err != nil {
		log.Fatal(err)
	} else {
		v.S()
	}

	if v, err := em.GetCarbonIntensityHistory("IN-EA", electricitymaps.OptionEmissionFactorType("direct")); err != nil {
		log.Fatal(err)
	} else {
		v.S()
	}
}

func handleEmber() {
	e := ember.DefaultEmber

	if v, err := e.GetCountries(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(v)
	}

	if v, err := e.GetMonthlyCo2Intensity("IND"); err != nil {
		log.Fatal(err)
	} else {
		v.S()
	}
}

func main() {
	// handleClimateTrace()
	handleElectricityMaps()
	// handleEmber()
}
