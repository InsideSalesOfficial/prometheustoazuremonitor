package main

import (
	"log"

	"github.com/InsideSalesOfficial/prometheustoazuremonitor/cfg"
	"github.com/InsideSalesOfficial/prometheustoazuremonitor/prometheus"
)

func main() {
	// Gather environment configuration and exit if we don't have all we need
	c, err := cfg.New()
	if err != nil {
		log.Fatal(err)
	}

	prometheus.PullPrometheusAndSend(c)
}
