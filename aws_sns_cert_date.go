package main

import (
	"errors"
	"time"

	"github.com/evalphobia/aws-sdk-go-wrapper/config"
	"github.com/evalphobia/aws-sdk-go-wrapper/sns"
)

var (
	ErrNoAppleCertificateExpiration = errors.New(`no "AppleCertificateExpirationDate" key in the attribute response`)
)

// GetAPNsCertDate fetchs the expiration date of apple cert on AWS SNS.
func GetAPNsCertDate(arn string) (time.Time, error) {
	cli, err := sns.New(config.Config{}, sns.Platforms{})
	if err != nil {
		return time.Time{}, err
	}

	attrs, err := cli.GetPlatformApplicationAttributes(arn)
	if err != nil {
		return time.Time{}, err
	}
	if !attrs.HasAppleCertificateExpirationDate {
		return time.Time{}, ErrNoAppleCertificateExpiration
	}
	return attrs.AppleCertificateExpirationDate, nil
}
