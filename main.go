package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordMetrics() {
	c := time.Tick(15 * time.Second)
	// for now := range c {
	for range c {
		log.Printf("Validating AWS data again...")

		data, _ := getAWSData()
		for _, r := range data {
			activeInstances.WithLabelValues(
				r.ID,
				r.InstanceType,
				r.Platform,
				r.OfferClass,
				r.OfferType,
				fmt.Sprintf("%v", r.Start),
				fmt.Sprintf("%d", r.Duration),
				fmt.Sprintf("%v", r.End),
				fmt.Sprintf("%.2f", r.TimeLeft),
			).Set(r.Count)
		}
	}
}

var (
	activeInstances = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "reserved",
		Subsystem: "active_instances",
		Name:      "count",
		Help:      "Number of active reserved instances.",
	}, []string{
		"RI_ID",
		"instance_type",
		"platform",
		"offer_class",
		"offer_type",
		"start",
		"duration",
		"end",
		"left",
	})
)

func main() {
	go recordMetrics()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":3001", nil)
}
