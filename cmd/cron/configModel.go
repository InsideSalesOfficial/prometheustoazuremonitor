// To parse and unparse this JSON data, add this code to your project and do:
//
//    config, err := UnmarshalConfig(bytes)
//    bytes, err = config.Marshal()

package main

import "encoding/json"

func UnmarshalConfig(data []byte) (Config, error) {
	var r Config
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Config) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Config struct {
	Config []ConfigElement `json:"config"`
}

type ConfigElement struct {
	PromURL               string   `json:"promURL"`
	AzureMonitorNamespace string   `json:"azureMonitorNamespace"`
	Metrics               []string `json:"metrics"`
}
