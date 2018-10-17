package main

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var (
	sess *session.Session
	svc  *ec2.EC2
)
