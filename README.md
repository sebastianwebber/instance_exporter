# instance_exporter

Prometheus exporter for ec2 instances metrics.

## Dependencies

Download AWS SDK:
```
go get -u -v github.com/aws/aws-sdk-go
```

Setup the follow environment variables:
- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`

## how to build (and run)

Use [gin](https://github.com/codegangsta/gin)...

```
gin
```

... And open [http://127.0.0.1:3000/metrics](http://127.0.0.1:3000/metrics)!