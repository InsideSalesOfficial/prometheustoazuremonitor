package cfg

import (
	"encoding/json"
	"io/ioutil"

	"github.com/kelseyhightower/envconfig"
)

// Cfg represents the possible values we can be configured with
// via environment variables and from a file.
type Cfg struct {
	AzureADClientID     string `json:"-" envconfig:"AZURE_AD_CLIENT_ID" required:"true"`
	AzureADTenantID     string `json:"-" envconfig:"AZURE_AD_TENANT_ID" required:"true"`
	AzureADClientSecret string `json:"-" envconfig:"AZURE_AD_CLIENT_SECRET" required:"true"`
	AzureResourceID     string `json:"-" envconfig:"AZURE_RESOURCE_ID" required:"true"`
	AzureMonitorRegion  string `json:"-" envconfig:"AZURE_MONITOR_REGION" default:"eastus"`
	// MetricsConfigFile is where we read in Config from
	MetricsConfigFile string          `json:"-" envconfig:"METRICS_CONFIG_FILE" default:"/mymnt/metrics.conf"`
	Config            []ConfigElement `json:"config"`
}

// New returns a new configuration, populated from environment variables and/or defaults.
// Also pulls in configuration from a configfile.
func New() (*Cfg, error) {
	cfg := &Cfg{}
	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, err
	}

	// Read and decode the metrics configuration JSON
	data, err := ioutil.ReadFile(cfg.MetricsConfigFile)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// ConfigElement represents each Prometheus-Azure Monitor Namespace set of metrics
type ConfigElement struct {
	PromURL               string   `json:"promURL"`
	AzureMonitorNamespace string   `json:"azureMonitorNamespace"`
	Metrics               []string `json:"metrics"`
}
