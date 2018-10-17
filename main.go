package main

import (
	"fmt"
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
	c := time.Tick(15 * time.Second)
	// for now := range c {
	for range c {
		log.Printf("Validating AWS data again...")

		data, _ := getReservedInstances()
		for _, r := range data {
			activeReservations.WithLabelValues(
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

func init() {
	sess = session.Must(session.NewSession(&aws.Config{
		Region: aws.String(endpoints.UsEast1RegionID),
	}))
	svc = ec2.New(sess)
}

func main() {
	go recordMetrics()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":3001", nil)
}
