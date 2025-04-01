package ember

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Ember struct{}

func NewEmber() *Ember {
	return &Ember{}
}

var DefaultEmber = NewEmber()

func fetchAuthToken() (string, error) {
	v, ok := os.LookupEnv("EMBER_TOKEN")
	if !ok || v == "" {
		return "", fmt.Errorf("missing environment variable EMBER_TOKEN")
	}

	return v, nil
}

func (c *Ember) GetCountries() ([]string, error) {
	authToken, err := fetchAuthToken()
	if err != nil {
		return nil, err
	}

	url := "https://api.ember-energy.org/v1/options/carbon-intensity/monthly/entity_code"

	url = fmt.Sprintf("%s?api_key=%s", url, authToken)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("API returned non-200 status code: %d, response: %s",
			res.StatusCode, string(bodyBytes))
	}

	type EntityRes struct {
		Options []string `json:"options"`
	}
	v := EntityRes{}

	if err := json.NewDecoder(res.Body).Decode(&v); err != nil {
		return nil, err
	}

	var countries []string
	for _, country := range v.Options {
		countries = append(countries, country)
	}

	return countries, nil
}

type Co2Intensity struct {
	Entity            string  `json:"entity"`
	EntityCode        string  `json:"entity_code"`
	IsAggregateEntity bool    `json:"is_aggregate_entity"`
	Date              string  `json:"date"`
	EmissionIntensity float64 `json:"emissions_intensity_gco2_per_kwh"`
	// Unit gco2_per_kwh
	Unit string `json:"unit"`
}

type MonthlyCo2Intensity []Co2Intensity

func (v MonthlyCo2Intensity) S() {
	for _, c := range v {
		fmt.Printf("%s (%s::%s) IsAggregateEntity: %t, Co2Intensity: %.2f%s\n",
			c.Date, c.Entity, c.EntityCode, c.IsAggregateEntity, c.EmissionIntensity, c.Unit)
	}
}

func (c *Ember) GetMonthlyCo2Intensity(countryCode string) (MonthlyCo2Intensity, error) {
	url := "https://api.ember-energy.org/v1/carbon-intensity/monthly"

	query := []string{"is_aggregate_entity=false", "include_all_dates_value_range=false", "entity_code=" + countryCode}

	authToken, err := fetchAuthToken()
	if err != nil {
		return nil, err
	}

	query = append(query, "api_key="+authToken)

	url += fmt.Sprintf("?%s", strings.Join(query, "&"))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("API returned non-200 status code: %d, response: %s",
			res.StatusCode, string(bodyBytes))
	}

	type EntityRes struct {
		Options MonthlyCo2Intensity `json:"data"`
	}
	v := EntityRes{}

	if err := json.NewDecoder(res.Body).Decode(&v); err != nil {
		return nil, err
	}

	for i := range v.Options {
		v.Options[i].Unit = "gCO2eq/kWh"
	}

	return v.Options, nil
}
