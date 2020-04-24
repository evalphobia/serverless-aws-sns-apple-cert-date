serverless-aws-sns-apple-cert-date
----

[![GoDoc][1]][2] [![License: MIT][3]][4] [![Release][5]][6] [![Build Status][7]][8] [![Codecov Coverage][11]][12] [![Go Report Card][13]][14] [![Code Climate][19]][20] [![BCH compliance][21]][22]

[1]: https://godoc.org/github.com/evalphobia/serverless-aws-sns-apple-cert-date?status.svg
[2]: https://godoc.org/github.com/evalphobia/serverless-aws-sns-apple-cert-date
[3]: https://img.shields.io/badge/License-MIT-blue.svg
[4]: LICENSE.md
[5]: https://img.shields.io/github/release/evalphobia/serverless-aws-sns-apple-cert-date.svg
[6]: https://github.com/evalphobia/serverless-aws-sns-apple-cert-date/releases/latest
[7]: https://github.com/evalphobia/serverless-aws-sns-apple-cert-date/workflows/test/badge.svg
[8]: https://github.com/evalphobia/serverless-aws-sns-apple-cert-date/actions?query=workflow%3Atest
[9]: https://coveralls.io/repos/evalphobia/serverless-aws-sns-apple-cert-date/badge.svg?branch=master&service=github
[10]: https://coveralls.io/github/evalphobia/serverless-aws-sns-apple-cert-date?branch=master
[11]: https://codecov.io/github/evalphobia/serverless-aws-sns-apple-cert-date/coverage.svg?branch=master
[12]: https://codecov.io/github/evalphobia/serverless-aws-sns-apple-cert-date?branch=master
[13]: https://goreportcard.com/badge/github.com/evalphobia/serverless-aws-sns-apple-cert-date
[14]: https://goreportcard.com/report/github.com/evalphobia/serverless-aws-sns-apple-cert-date
[15]: https://img.shields.io/github/downloads/evalphobia/serverless-aws-sns-apple-cert-date/total.svg?maxAge=1800
[16]: https://github.com/evalphobia/serverless-aws-sns-apple-cert-date/releases
[17]: https://img.shields.io/github/stars/evalphobia/serverless-aws-sns-apple-cert-date.svg
[18]: https://github.com/evalphobia/serverless-aws-sns-apple-cert-date/stargazers
[19]: https://codeclimate.com/github/evalphobia/serverless-aws-sns-apple-cert-date/badges/gpa.svg
[20]: https://codeclimate.com/github/evalphobia/serverless-aws-sns-apple-cert-date
[21]: https://bettercodehub.com/edge/badge/evalphobia/serverless-aws-sns-apple-cert-date?branch=master
[22]: https://bettercodehub.com/

`serverless-aws-sns-apple-cert-date` checks AppleCertificateExpirationDate on AWS SNS and send slack when expiring, powered by AWS Lambda.


# Download

Download serverless-aws-sns-apple-cert-date by command below.

```bash
$ git clone https://github.com/evalphobia/serverless-aws-sns-apple-cert-date
$ cd serverless-aws-sns-apple-cert-date
$ make init
```

# Config

## serverless.yml

Change environment variables below,

```bash
$ vim serverless.yml

------------

provider:
  name: aws
  region: ap-northeast-1  # <- Change to your target region.


...
functions:
  check:
    handler: bin/serverless
    memorySize: 128
    timeout: 119
    environment:
      # Change to your target Application Platform ARN of AWS SNS.
      ARN_LIST: >-
        arn:aws:sns:ap-northeast-1:000000000000:app/APNS/app1
        arn:aws:sns:ap-northeast-1:000000000000:app/APNS/app2
        arn:aws:sns:ap-northeast-1:000000000000:app/APNS/app3
        arn:aws:sns:ap-northeast-1:000000000000:app/APNS/app4
      # Change to your own threshold.
      TIME_THRESHOLD: 30d10h # 30days and 10 hours before
      # If you want to get report even if it's not in expiring, turn it to true.
      FORCE_REPORT: false

      # SLACK_WEBHOOK_URL_KMS: 'xxx'
      SLACK_WEBHOOK_URL: https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX
      SLACK_CHANNEL: '#channel'
    events:
      - schedule: cron(0 0 * * ? *)  # exec everyday on 00:00
```

## Environment variables

|Name|Description|Default|
|:--|:--|:--|
| `ARN_LIST` | Puts Application Platform ARNs on AWS SNS. | - |
| `FORCE_REPORT` | A flag to send the report to Slack even if not in expiring. | `false` |
| `TIME_THRESHOLD` | Expiration time threshold to send report. Besides go's ParseDuration format, it supports `d` for days. (ref: https://golang.org/pkg/time/#ParseDuration ) | `1d` |
| `SLACK_WEBHOOK_URL` | Slack's webhook URL. | - |
| `SLACK_WEBHOOK_URL_KMS` | Slack's webhook URL encrypted by AWS KMS. | `false` |
| `SLACK_USERNAME` | A sender name on Slack. | - |
| `SLACK_CHANNEL` | A channel name on Slack. | - |


# Deploy

```bash
$ AWS_ACCESS_KEY_ID=<...> AWS_SECRET_ACCESS_KEY=<...> make deploy
```


# Check Log

```bash
$ AWS_ACCESS_KEY_ID=<...> AWS_SECRET_ACCESS_KEY=<...> sls logs -f <function name> -t
```
