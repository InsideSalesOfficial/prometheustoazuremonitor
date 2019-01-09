// To parse and unparse this JSON data, add this code to your project and do:
//
//    AzureMonitor, err := UnmarshalAzureMonitor(bytes)
//    bytes, err = AzureMonitor.Marshal()

package azuremonitor

import (
	"encoding/json"
	"time"
)

//UnmarshalAzureMonitor parses the json to the AzureMonitor
func UnmarshalAzureMonitor(data []byte) (AzureMonitor, error) {
	var r AzureMonitor
	err := json.Unmarshal(data, &r)
	return r, err
}

//Marshal parses to Json
func (r *AzureMonitor) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

//AzureMonitor AzureMonitor Azure Monitor Json
type AzureMonitor struct {
	Timestamp time.Time `json:"time"`
	Data      Data      `json:"data"`
}

//Data define the data structure to be sent to AzureMonitor
type Data struct {
	BaseData BaseData `json:"baseData"`
}

//BaseData Define the base data
type BaseData struct {
	Metric    string   `json:"metric"`
	Namespace string   `json:"namespace"`
	DimNames  []string `json:"dimNames"`
	Series    []Series `json:"series"`
}

//Series Define series
type Series struct {
	DimValues []string `json:"dimValues"`
	Min       float64  `json:"min"`
	Max       float64  `json:"max"`
	Sum       float64  `json:"sum"`
	Count     int64    `json:"count"`
}
