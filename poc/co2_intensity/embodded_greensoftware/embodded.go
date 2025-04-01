package embodded_greensoftware

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"net/http"
	"strconv"
	"strings"
)

// VMEmission represents CO2 emission data for a specific VM type
type VMEmission struct {
	Series       string  `json:"series"`
	InstanceType string  `json:"instance_type"`
	EmboddedCo2  float64 `json:"embodded_co2"`
	Co2Unit      string  `json:"co2_unit"`
}

type EmboddedCo2Emissions struct {
	Azure map[string]VMEmission `json:"azure"`
	GCP   map[string]VMEmission `json:"gcp"`
	AWS   map[string]VMEmission `json:"aws"`
}

func (e *EmboddedCo2Emissions) P() {
	b, _ := json.MarshalIndent(e, "", " ")
	fmt.Println(string(b))
}

func (e *EmboddedCo2Emissions) S() {
	color.HiCyan("Azure:")
	for k, v := range e.Azure {
		fmt.Printf("%s: %f%s\n", k, v.EmboddedCo2, v.Co2Unit)
	}

	color.HiCyan("GCP:")
	for k, v := range e.GCP {
		fmt.Printf("%s: %f%s\n", k, v.EmboddedCo2, v.Co2Unit)
	}

	color.HiCyan("AWS:")
	for k, v := range e.AWS {
		fmt.Printf("%s: %f%s\n", k, v.EmboddedCo2, v.Co2Unit)
	}
}

// GetEmboddedCo2Emissions fetches and parses embodied CO2 emissions for cloud providers
func GetEmboddedCo2Emissions() (*EmboddedCo2Emissions, error) {
	url := "https://raw.githubusercontent.com/ksctl/components/refs/heads/main/co2/embodied_emissions.csv"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer res.Body.Close()

	reader := csv.NewReader(res.Body)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %v", err)
	}

	records = records[1:]

	type coordinates struct {
		start, end int
	}

	azureRecords := coordinates{0, 6}
	gcpRecords := coordinates{7, 13}
	awsRecords := coordinates{14, 20}

	dataAws := make(map[string]VMEmission, 200)
	dataAzure := make(map[string]VMEmission, 200)
	dataGcp := make(map[string]VMEmission, 100)

	for _, record := range records {
		aws := record[awsRecords.start:awsRecords.end]
		azure := record[azureRecords.start:azureRecords.end]
		gcp := record[gcpRecords.start:gcpRecords.end]

		if len(gcp[1]) > 0 {
			gcpEmission, _ := strconv.ParseFloat(gcp[4], 64)
			dataGcp[gcp[1]] = VMEmission{
				Series:       gcp[0],
				InstanceType: gcp[1],
				EmboddedCo2:  gcpEmission,
				Co2Unit:      "kgCO₂eq",
			}
		}

		if len(azure[1]) > 0 && (strings.HasPrefix(azure[1], "E") ||
			strings.HasPrefix(azure[1], "D") ||
			strings.HasPrefix(azure[1], "B") ||
			strings.HasPrefix(azure[1], "F")) {
			azureEmission, _ := strconv.ParseFloat(azure[4], 64)
			sku := "Standard_" + strings.ReplaceAll(azure[1], " ", "_")
			dataAzure[sku] = VMEmission{
				Series:       azure[0],
				InstanceType: sku,
				EmboddedCo2:  azureEmission,
				Co2Unit:      "kgCO₂eq",
			}
		}

		if len(aws[1]) > 0 && (strings.HasPrefix(aws[1], "m5") ||
			strings.HasPrefix(aws[1], "c5") ||
			strings.HasPrefix(aws[1], "r5") ||
			strings.HasPrefix(aws[1], "t3")) {
			awsEmission, _ := strconv.ParseFloat(aws[4], 64)
			dataAws[aws[1]] = VMEmission{
				Series:       aws[0],
				InstanceType: aws[1],
				EmboddedCo2:  awsEmission,
				Co2Unit:      "kgCO₂eq",
			}
		}
	}

	finalData := EmboddedCo2Emissions{
		Azure: dataAzure,
		GCP:   dataGcp,
		AWS:   dataAws,
	}

	return &finalData, nil
}
