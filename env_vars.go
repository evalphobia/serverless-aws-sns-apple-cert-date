package main

import (
	"os"
	"strings"
)

func GetEnvARNList() string {
	return os.Getenv("ARN_LIST")
}

func GetEnvForceReport() string {
	return os.Getenv("FORCE_REPORT")
}

func GetEnvTimeThreshold() string {
	return os.Getenv("TIME_THRESHOLD")
}

func GetEnvSlackUser() string {
	return os.Getenv("SLACK_USERNAME")
}

func GetEnvSlackChannel() string {
	return os.Getenv("SLACK_CHANNEL")
}

func GetEnvSlackURL() string {
	kmsText := os.Getenv("SLACK_WEBHOOK_URL_KMS")
	if kmsText != "" {
		url, err := DecryptKMS(kmsText)
		if err == nil {
			return strings.TrimSpace(url)
		}
	}

	return os.Getenv("SLACK_WEBHOOK_URL")
}
