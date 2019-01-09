package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"

	"github.com/InsideSalesOfficial/prometheustoazuremonitor/azuremonitor"
	"github.com/InsideSalesOfficial/prometheustoazuremonitor/cfg"
)

func main() {
	// Gather environment configuration and exit if we don't have all we need
	c, err := cfg.New()
	if err != nil {
		panic(err)
	}

	// Decode the metrics configuration JSON
	dat, err := ioutil.ReadFile(c.MetricsConfigFile)
	if err != nil {
		panic(err)
	}
	config, err := UnmarshalConfig(dat)

	pullPrometheus(c, config)

}

func pullPrometheus(c *cfg.Cfg, config Config) {

	for _, promConfig := range config.Config {
		log.Println(promConfig)

		azureMonitorNamespace := promConfig.AzureMonitorNamespace
		client, err := api.NewClient(api.Config{
			Address: promConfig.PromURL,
		})

		if err != nil {
			panic(err)
		}

		promRange := v1.Range{End: time.Now(), Start: time.Now().Add(-time.Minute), Step: time.Minute}
		promAPI := v1.NewAPI(client)

		var wg sync.WaitGroup

		for _, metricName := range promConfig.Metrics {
			wg.Add(1)
			go func(metricName string) {
				promModel, err := promAPI.QueryRange(context.Background(), metricName, promRange)
				if err != nil {
					fmt.Println(err)
				}

				var promMatrix model.Matrix
				switch promModel.Type() {
				case model.ValMatrix:
					promMatrix = promModel.(model.Matrix)
				}

				if promMatrix == nil {
					return
				}

				var dimNames []string
				var series []azuremonitor.Series
				timestamp := time.Now()

				for _, p := range promMatrix {

					var dimValues []string
					m := make(model.Metric, len(p.Metric))
					for key, dim := range p.Metric {
						if key == "__name__" {
							continue
						}
						m[model.LabelName(key)] = model.LabelValue(dim)
					}

					//Creates a commom dimension name for all the series
					if len(series) == 0 {
						for mkey := range m {
							dimNames = append(dimNames, fmt.Sprintf("%s", mkey))
						}
					}

					for _, mname := range dimNames {
						dimValues = append(dimValues, fmt.Sprintf("%s", m[model.LabelName(mname)]))
					}

					var oldValue, currValue float64
					for _, metricValue := range p.Values {
						timestamp = metricValue.Timestamp.Time()
						if oldValue == 0 {
							oldValue = float64(metricValue.Value)
						} else {
							currValue = float64(metricValue.Value) - oldValue
						}
						// send zero if the monotically increasing counter decreased or reset during the time period
						if currValue < 0 {
							currValue = 0
						}
					}

					serie := azuremonitor.Series{DimValues: dimValues, Max: currValue, Min: currValue, Sum: currValue, Count: 1}
					series = append(series, serie)

				}

				var baseData azuremonitor.BaseData

				if len(series) == 0 {
					serie := azuremonitor.Series{Max: 0, Min: 0, Sum: 0, Count: 1}
					series = append(series, serie)
					baseData = azuremonitor.BaseData{Metric: metricName, Namespace: azureMonitorNamespace, Series: series}
				} else {
					baseData = azuremonitor.BaseData{Metric: metricName, Namespace: azureMonitorNamespace, DimNames: dimNames, Series: series}
				}

				data := azuremonitor.Data{BaseData: baseData}
				customData := azuremonitor.AzureMonitor{Timestamp: timestamp, Data: data}
				customDataBytes, err := customData.Marshal()
				if err != nil {
					log.Print(err)
					return
				}

				err = sendToAzureMonitor(c, string(customDataBytes))
				if err != nil {
					log.Print(err)
				}
			}(metricName)
		}
		wg.Wait()
	}

}

func sendToAzureMonitor(c *cfg.Cfg, postData string) error {
	var cli = azuremonitor.New(c.AzureADTenantID, c.AzureADClientID, c.AzureADClientSecret)
	fmt.Println(fmt.Sprintf("region: %s \n resourceID: %s \n postData: %s", c.AzureMonitorRegion, c.AzureResourceID, postData))
	err := cli.SaveCustomAzureData(c.AzureMonitorRegion, c.AzureResourceID, postData)
	if err != nil {
		log.Print(err)
	}
	return err
}
