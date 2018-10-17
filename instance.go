package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Instance contains a definition of a AWS EC2 Instance
type Instance struct {
	ID        string
	Type      string
	PublicDNS string
	NameTag   string
}

var (
	runningInstances = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "ec2",
		Subsystem: "running_instances",
		Name:      "count",
		Help:      "Number of active reserved instances.",
	}, []string{
		"id",
		"instance_type",
		"dns",
	})
)

func updateInstances() {
	log.Println("Update Instances data...")
	data, err := getEC2Instances()

	if err != nil {
		log.Fatalf("Could not get EC2 Instances: %v\n", err)
	}
	for _, r := range data {
		runningInstances.WithLabelValues(
			r.ID,
			r.Type,
			r.PublicDNS,
		).Set(1.0)

		fmt.Println(r.ID)
	}
}

func getEC2Instances() (output []Instance, err error) {

	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("instance-state-name"),
				Values: []*string{
					aws.String("running"),
				},
			},
		},
	}

	result, err := svc.DescribeInstances(input)
	if err != nil {
		return
	}

	for reserv := 0; reserv < len(result.Reservations); reserv++ {

		for i := 0; i < len(result.Reservations[reserv].Instances); i++ {

			fmt.Println(*result.Reservations[reserv].Instances[i])
			output = append(output, Instance{
				ID:        *result.Reservations[reserv].Instances[i].InstanceId,
				Type:      *result.Reservations[reserv].Instances[i].InstanceType,
				PublicDNS: *result.Reservations[reserv].Instances[i].PublicDnsName,
				// NameTag: *result.Reservations[reserv].Instances[i].Tags
			})
		}
	}
	return
}
