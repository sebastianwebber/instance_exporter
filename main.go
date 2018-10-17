package main

import (
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordMetrics() {
	log.Println("Starting instance_exporter...")

	updateReserved()
	reservedTick := time.NewTicker(30 * time.Minute)
	// for now := range c {
	go func() {
		for range reservedTick.C {
			updateReserved()
		}
	}()
}

func init() {
	sess = session.Must(session.NewSession(&aws.Config{
		Region: aws.String(endpoints.UsEast1RegionID),
	}))
	svc = ec2.New(sess)
}

func main() {
	go recordMetrics()

	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":3001", nil); err != nil {
		log.Fatalf("Could not start webserver: %v", err)
	}
}
