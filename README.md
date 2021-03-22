# Freedom Countries watcher

Watch the list of countries from
[FreedomHouse World](https://freedomhouse.org/countries/freedom-world/scores)
and
[FreedomHouse Net](https://freedomhouse.org/countries/freedom-net/scores) and
compute risk scores. The diff is sent via email.

## Deployment

Infra is maintained with Terraform. The main config is stored on S3.
Simply run `terraform init` the first time.

Terraform's variables are the following:

```
fcup_email = "xxx"
fcup_name  = "xxx"
s3_bucket  = "xxxx"
cron_rate  = "7 days"
```

[Cron rate
expression](https://docs.aws.amazon.com/AmazonCloudWatch/latest/events/ScheduledEvents.html#RateExpressions)

