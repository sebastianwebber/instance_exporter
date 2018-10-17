package main

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Reservation contains a definition of a AWS Instance Reservation
type Reservation struct {
	ID           string
	InstanceType string
	Platform     string
	OfferClass   string
	OfferType    string
	Start        time.Time
	End          time.Time
	Duration     int64
	TimeLeft     float64
	Count        float64
}

var (
	activeReservations = promauto.NewGaugeVec(prometheus.GaugeOpts{
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

func getReservedInstances() (output []Reservation, err error) {

	input := &ec2.DescribeReservedInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("state"),
				Values: []*string{
					aws.String("active"),
				},
			},
		},
	}

	result, err := svc.DescribeReservedInstances(input)
	if err != nil {
		return
	}

	for i := 0; i < len(result.ReservedInstances); i++ {

		duration := *result.ReservedInstances[i].Duration
		startDate := *result.ReservedInstances[i].Start
		endDate := startDate.Add(time.Second * time.Duration(duration))
		left := endDate.Sub(startDate)

		output = append(output, Reservation{
			ID:           *result.ReservedInstances[i].ReservedInstancesId,
			InstanceType: *result.ReservedInstances[i].InstanceType,
			Platform:     *result.ReservedInstances[i].ProductDescription,
			OfferClass:   *result.ReservedInstances[i].OfferingClass,
			OfferType:    *result.ReservedInstances[i].OfferingType,
			Count:        float64(*result.ReservedInstances[i].InstanceCount),
			Start:        startDate,
			End:          endDate,
			Duration:     duration,
			TimeLeft:     left.Seconds(),
		})
	}
	return
}