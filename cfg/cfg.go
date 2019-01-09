package cfg

import "github.com/kelseyhightower/envconfig"

// Cfg represents the possible values edgy can be configured with via environment variables.
/*
var configurationFile = os.Getenv("METRICS_CONFIG_FILE")
var clientID = os.Getenv("AZURE_AD_CLIENT_ID")
var tenantID = os.Getenv("AZURE_AD_TENANT_ID")
var clientSecret = os.Getenv("AZURE_AD_CLIENT_SECRET")
var resourceID = os.Getenv("AZURE_RESOURCE_ID")
var region = os.Getenv("AZURE_MONITOR_REGION")
*/
type Cfg struct {
	MetricsConfigFile   string `envconfig:"METRICS_CONFIG_FILE" default:"/mymnt/metrics.conf"`
	AzureADClientID     string `envconfig:"AZURE_AD_CLIENT_ID" required:"true"`
	AzureADTenantID     string `envconfig:"AZURE_AD_TENANT_ID" required:"true"`
	AzureADClientSecret string `envconfig:"AZURE_AD_CLIENT_SECRET" required:"true"`
	AzureResourceID     string `envconfig:"AZURE_RESOURCE_ID" required:"true"`
	AzureMonitorRegion  string `envconfig:"AZURE_MONITOR_REGION" default:"eastus"`
}

// New returns a new configuration, populated from environment variables and/or defaults.
func New() (*Cfg, error) {
	cfg := &Cfg{}
	return cfg, envconfig.Process("", cfg)
}
