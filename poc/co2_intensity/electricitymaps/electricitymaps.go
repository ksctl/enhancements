package electricitymaps

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

type ElectricityMap struct{}

func NewElectricityMap() *ElectricityMap {
	return &ElectricityMap{}
}

var DefaultElectricityMap = NewElectricityMap()

func fetchAuthToken() (string, error) {
	authToken, ok := os.LookupEnv("ELECTRICITY_MAP_TOKEN")
	if !ok || authToken == "" {
		return "", fmt.Errorf("missing environment variable ELECTRICITY_MAP_TOKEN")
	}

	return authToken, nil
}

type ZoneInfo struct {
	ZoneName string   `json:"zoneName"`
	Access   []string `json:"access"`
}

type ZonesResponse map[string]ZoneInfo

func (v ZonesResponse) P() {
	b, _ := json.MarshalIndent(v, "", " ")

	fmt.Println(string(b))
}

func (v ZonesResponse) S() {
	color.HiCyan("Zones:")
	for k, v := range v {
		fmt.Println(k, "->", v.ZoneName)
	}
}

func (c *ElectricityMap) GetAvailableZones() (ZonesResponse, error) {
	authToken, err := fetchAuthToken()
	if err != nil {
		return nil, err
	}

	url := "https://api.electricitymap.org/v3/zones"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("auth-token", authToken)

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

	var zones ZonesResponse
	if err := json.NewDecoder(res.Body).Decode(&zones); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %v", err)
	}

	return zones, nil
}

type CarbonIntensityRecord struct {
	Year     int
	Month    string
	Country  string
	ZoneName string
	ZoneID   string
	// DirectCarbonIntensity has unit gCO2eq/kWh
	DirectCarbonIntensity float64
	// LCACarbonIntensity has unit gCO2eq/kWh
	LCACarbonIntensity  float64
	LowCarbonPercentage float64
	RenewablePercentage float64
	DataSource          string
	Unit                string
}

type CarbonIntensityRecords []CarbonIntensityRecord

func (v CarbonIntensityRecords) P() {
	b, _ := json.MarshalIndent(v, "", " ")

	color.HiCyan("Carbon Intensity Records:")
	fmt.Println(string(b))
}

func (v CarbonIntensityRecords) S() {
	color.HiCyan("Carbon Intensity Records:")
	for _, r := range v {
		fmt.Printf("%d-%s, (%s::%s), DirectCo2Intensity: %f, LCACo2Intensity: %f, LowCarbonPercentage: %f, RenewablePercentage: %f, Unit: %s\n",
			r.Year, r.Month, r.Country, r.ZoneID, r.DirectCarbonIntensity, r.LCACarbonIntensity, r.LowCarbonPercentage, r.RenewablePercentage, r.Unit)
	}
}

func (c *ElectricityMap) GetMonthlyPastData(zoneId string) (CarbonIntensityRecords, error) {
	url := fmt.Sprintf("https://data.electricitymaps.com/2025-01-27/%s_2024_monthly.csv", zoneId)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	reader := csv.NewReader(res.Body)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("no data found in CSV")
	}

	var data CarbonIntensityRecords
	for i, row := range records[1:] { // Skipping header row
		if len(row) < 9 {
			return nil, fmt.Errorf("unexpected number of columns in row %d", i+1)
		}

		timestamp, err := time.Parse("2006-01-02 15:04:05", row[0])
		if err != nil {
			return nil, fmt.Errorf("error parsing date on row %d: %v", i+1, err)
		}

		directCI, err := strconv.ParseFloat(row[4], 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing direct carbon intensity on row %d: %v", i+1, err)
		}

		lcaCI, err := strconv.ParseFloat(row[5], 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing LCA carbon intensity on row %d: %v", i+1, err)
		}

		lowCarbonPercent, err := strconv.ParseFloat(row[6], 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing low carbon percentage on row %d: %v", i+1, err)
		}

		renewablePercent, err := strconv.ParseFloat(row[7], 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing renewable percentage on row %d: %v", i+1, err)
		}

		data = append(data, CarbonIntensityRecord{
			Year:                  timestamp.Year(),
			Month:                 timestamp.Month().String(),
			Country:               row[1],
			ZoneName:              row[2],
			ZoneID:                row[3],
			DirectCarbonIntensity: directCI,
			LCACarbonIntensity:    lcaCI,
			LowCarbonPercentage:   lowCarbonPercent,
			RenewablePercentage:   renewablePercent,
			DataSource:            row[8],
			Unit:                  "gCO2eq/kWh",
		})
	}

	return data, nil
}

type option struct {
	emissionFactorType string
	disableEstimations bool
}

type Option func(*option) error

func OptionEmissionFactorType(emissionFactorType string) Option {
	return func(o *option) error {
		// lifecycle includes all the emissions from the production of electricity (e.g. coal mining, gas extraction, etc.)
		// Direct emissions are typically used for immediate operational assessments, such as calculating Scope 1 and Scope 2 emissions
		if emissionFactorType != "lifecycle" && emissionFactorType != "direct" {
			return fmt.Errorf("invalid emission factor type")
		}
		o.emissionFactorType = emissionFactorType
		return nil
	}
}

func OptionDisableEstimations() Option {
	return func(o *option) error {
		o.disableEstimations = true
		return nil
	}
}

type CarbonIntensity struct {
	Co2Intensity       float64   `json:"carbonIntensity"`
	Datetime           time.Time `json:"datetime"`
	UpdatedAt          time.Time `json:"updatedAt"`
	EmissionFactorType string    `json:"emissionFactorType"`
	IsEstimated        bool      `json:"isEstimated"`
	EstimationMethod   *string   `json:"estimationMethod"`
	Unit               string
}

type CarbonIntensityLatest struct {
	Zone string `json:"zone"`
	CarbonIntensity
}

func (v *CarbonIntensityLatest) P() {
	b, _ := json.MarshalIndent(v, "", " ")

	color.HiCyan("Carbon Intensity Latest:")
	fmt.Println(string(b))
}

func (v *CarbonIntensityLatest) S() {
	color.HiCyan("Carbon Intensity Latest:")
	fmt.Printf("%s [%s] CarbonIntensity: %f%s, {%s}\n", v.UpdatedAt, v.Zone, v.Co2Intensity, v.Unit, v.Datetime)
}

func (c *ElectricityMap) GetLatestCarbonIntensity(zoneID string, opts ...Option) (*CarbonIntensityLatest, error) {
	authToken, err := fetchAuthToken()
	if err != nil {
		return nil, err
	}

	var o option
	for _, opt := range opts {
		if err := opt(&o); err != nil {
			return nil, err
		}
	}

	url := "https://api.electricitymap.org/v3/carbon-intensity/latest"
	query := []string{}
	if zoneID != "" {
		query = append(query, "zone="+zoneID)
	}

	if o.emissionFactorType != "" {
		query = append(query, "emissionFactorType="+o.emissionFactorType)
	}
	if o.disableEstimations {
		query = append(query, "disableEstimations=true")
	}

	if len(query) > 0 {
		url += "?" + strings.Join(query, "&")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("auth-token", authToken)

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

	var data CarbonIntensityLatest
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %v", err)
	}

	data.Unit = "gCO2eq/kWh"

	return &data, nil
}

type CarbonIntensityHis struct {
	CreatedAt time.Time `json:"createdAt"`
	CarbonIntensity
}

type CarbonIntensityHistory struct {
	Zone    string               `json:"zone"`
	History []CarbonIntensityHis `json:"history"`
}

func (v *CarbonIntensityHistory) P() {
	b, _ := json.MarshalIndent(v, "", " ")

	color.HiCyan("Carbon Intensity History(24h):")
	fmt.Println(string(b))
}

func (v *CarbonIntensityHistory) S() {
	color.HiCyan("Carbon Intensity History(24h):")
	for _, h := range v.History {
		fmt.Printf("%s [%s] CarbonIntensity: %f%s, {%s}\n", h.CreatedAt, v.Zone, h.Co2Intensity, h.Unit, h.UpdatedAt)
	}
}

func (c *ElectricityMap) GetCarbonIntensityHistory(zoneID string, opts ...Option) (*CarbonIntensityHistory, error) {
	authToken, err := fetchAuthToken()
	if err != nil {
		return nil, err
	}

	var o option
	for _, opt := range opts {
		if err := opt(&o); err != nil {
			return nil, err
		}
	}

	url := "https://api.electricitymap.org/v3/carbon-intensity/history"

	query := []string{}
	if zoneID != "" {
		query = append(query, "zone="+zoneID)
	}

	if o.emissionFactorType != "" {
		query = append(query, "emissionFactorType="+o.emissionFactorType)
	}
	if o.disableEstimations {
		query = append(query, "disableEstimations=true")
	}

	if len(query) > 0 {
		url += "?" + strings.Join(query, "&")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("auth-token", authToken)

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

	var data CarbonIntensityHistory
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %v", err)
	}

	for i := range data.History {
		data.History[i].Unit = "gCO2eq/kWh"
	}

	return &data, nil
}
