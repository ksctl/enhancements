package climatetrace

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ClimateTrace struct{}

var DefaultClimateTrace = NewClimateTrace()

func NewClimateTrace() *ClimateTrace {
	return &ClimateTrace{}
}

type Country struct {
	Name      string `json:"name"`
	Continent string `json:"continent"`
	Code      string `json:"alpha3"`
}

type Countries []Country

func (v Countries) P() {
	b, _ := json.MarshalIndent(v, "", " ")

	fmt.Println(string(b))
}

func (c *ClimateTrace) GetCountries() (Countries, error) {
	url := "https://api.climatetrace.org/v6/definitions/countries"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	countries := Countries{}
	if err := json.NewDecoder(res.Body).Decode(&countries); err != nil {
		return nil, err
	}

	return countries, nil
}

type Emissions struct {
	CO2       float64 `json:"co2"`
	CH4       float64 `json:"ch4"`
	N2O       float64 `json:"n2o"`
	CO2E100yr float64 `json:"co2e_100yr"`
	CO2E20yr  float64 `json:"co2e_20yr"`
}

type CountryEmission struct {
	Country      string    `json:"country"`
	Continent    *string   `json:"continent"`
	Rank         int       `json:"rank"`
	PreviousRank int       `json:"previousRank"`
	AssetCount   *int      `json:"assetCount"`
	Emissions    Emissions `json:"emissions"`
	// WorldEmissions  Emissions `json:"worldEmissions"`
	EmissionsChange Emissions `json:"emissionsChange"`
}

type CountriesEmission []CountryEmission

func (v CountriesEmission) P() {
	b, _ := json.MarshalIndent(v, "", " ")

	fmt.Println(string(b))
}

func (c *ClimateTrace) GetEmissionSummaryHistroy() (CountriesEmission, error) {
	url := "https://api.climatetrace.org/v6/country/emissions?since=2023&to=2024&countries=IND,FIN"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	countriesEmission := CountriesEmission{}
	if err := json.NewDecoder(res.Body).Decode(&countriesEmission); err != nil {
		return nil, err
	}

	return countriesEmission, nil
}
