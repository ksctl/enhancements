package main

import (
	"log"

	"github.com/ksctl/enhancements/poc/co2_intensity/climatetrace"
	"github.com/ksctl/enhancements/poc/co2_intensity/electricitymaps"
)

func climateTrace() {
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

func electricityMaps() {
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

func main() {
	// climateTrace()
	electricityMaps()
}
