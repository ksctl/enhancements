package climatiq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/fatih/color"
)

type Climatiq struct{}

func NewClimatiq() *Climatiq {
	return &Climatiq{}
}

var DefaultClimatiq = NewClimatiq()

func fetchAuthToken() (string, error) {
	v, ok := os.LookupEnv("CLIMATIQ_API_KEY")
	if !ok || v == "" {
		return "", fmt.Errorf("missing environment variable CLIMATIQ_API_KEY")
	}

	return v, nil
}

type CloudProvider struct {
	ProviderFullName        string   `json:"provider_full_name"`
	ProviderID              string   `json:"provider_id"`
	Regions                 []string `json:"regions"`
	VirtualMachineInstances []string `json:"virtual_machine_instances,omitempty"`
}

type CloudProvidersData struct {
	CloudProviders map[string]CloudProvider `json:"cloud_providers"`
}

func (v CloudProvidersData) P() {
	color.HiCyan("Cloud Providers:")
	b, _ := json.MarshalIndent(v, "", " ")

	fmt.Println(string(b))
}

func (c *Climatiq) GetMetadata() (*CloudProvidersData, error) {
	url := "https://api.climatiq.io/compute/v1/metadata"

	authToken, err := fetchAuthToken()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authToken))

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

	var data CloudProvidersData

	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %v", err)
	}

	return &data, nil
}

// GetDataBasedProvider It required paid plan so left this provider
func (c *Climatiq) GetDataBasedProvider(provider string) error {
	//if provider != "aws" || provider != "gcp" || provider != "azure" {
	//	return fmt.Errorf("invalid provider: %s", provider)
	//}
	provider = "aws"
	url := fmt.Sprintf("https://api.climatiq.io/compute/v1/%s/instance", provider)

	authToken, err := fetchAuthToken()
	if err != nil {
		return err
	}

	payloadBuf := new(bytes.Buffer)

	data := map[string]any{
		"region":        "us-east-1",
		"instance":      "c5.large",
		"duration":      730,
		"duration_unit": "h",
	}

	if err := json.NewEncoder(payloadBuf).Encode(data); err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, payloadBuf)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authToken))
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		return fmt.Errorf("API returned non-200 status code: %d, response: %s",
			res.StatusCode, string(bodyBytes))
	}

	bodyBytes, _ := io.ReadAll(res.Body)
	fmt.Println(string(bodyBytes))

	return nil
}
