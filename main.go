package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context, event LambdaEvent) (LambdaResponse, error) {
	list := event.getARNList()
	switch {
	case len(list) == 0:
		err := errors.New("ARN List are empty")
		errorf("err: [%v]", err)
		return NewErrResponse(err), err
	}

	timeThreshold, err := event.getTimeThreshold()
	if err != nil {
		err := errors.New("invalid TimeThreshold format")
		errorf("err: [%v]", err)
		return NewErrResponse(err), err
	}

	isForceReport := event.getForceReport()
	infof("list: [%v], threshold: [%v], force_report: [%v]", list, timeThreshold, isForceReport)

	// check ARNs
	var reportText []string
	for _, arn := range list {
		dt, err := GetAPNsCertDate(arn)
		if err != nil {
			errorf("GetAPNsCertDate err: [%v]", err)
			reportText = append(reportText, fmt.Sprintf("[%s]\t%v", arn, err))
		}
		if isForceReport || dt.Before(timeThreshold) {
			reportText = append(reportText, fmt.Sprintf("[%s]\t%v", arn, dt))
		}
		infof("ExpiresDate: [%s]\t%v", arn, dt)
	}

	// send report to slack
	if len(reportText) != 0 {
		err := sendReportToSlack(reportText)
		if err != nil {
			errorf("sendReportToSlack err: [%v]", err)
			return NewErrResponse(err), err
		}
	}

	return LambdaResponse{
		Status:  200,
		Success: true,
	}, nil
}

func sendReportToSlack(textList []string) error {
	msg := strings.Join(textList, "\n")

	req := SlackRequest{}
	attachment := SlackAttachment{
		Color:      "#ff9900",
		Title:      "SNS APNs Cert in Expiring",
		TitleLink:  "https://console.aws.amazon.com/sns/v3/home?#/mobile/push-notifications",
		MarkdownIn: []string{"text"},
		Text:       fmt.Sprintf("```%s```", msg),
	}
	req.Attachments = append(req.Attachments, attachment)
	return sendToSlack(req)
}

type LambdaResponse struct {
	Status  int   `json:"status"`
	Success bool  `json:"success"`
	Err     error `json:"error"`
}

func NewErrResponse(err error) LambdaResponse {
	return LambdaResponse{
		Status:  500,
		Success: false,
		Err:     err,
	}
}

type LambdaEvent struct {
	ARNList       string `json:"event_name"`
	ForceReport   string `json:"force_report"`
	TimeThreshold string `json:"time_threshold"`
}

func (e LambdaEvent) getARNList() []string {
	v := e.ARNList
	if v == "" {
		v = GetEnvARNList()
	}

	var list []string
	for _, vv := range strings.Split(strings.TrimSpace(v), " ") {
		vv = strings.TrimSpace(vv)
		if vv == "" {
			continue
		}
		list = append(list, vv)
	}
	return list
}

func (e LambdaEvent) getForceReport() bool {
	v := e.ForceReport
	if v == "" {
		v = GetEnvForceReport()
	}
	b, _ := strconv.ParseBool(v)
	return b
}

func (e LambdaEvent) getTimeThreshold() (time.Time, error) {
	v := e.TimeThreshold
	if v == "" {
		v = GetEnvTimeThreshold()
	}
	if v == "" {
		return time.Now().AddDate(0, 0, 1), nil
	}

	dt, err := time.Parse(time.RFC3339, v)
	if err == nil {
		return dt, nil
	}
	dur, err := parseDuration(v)
	if err != nil {
		return dt, err
	}
	return time.Now().Add(dur).UTC(), nil
}

func parseDuration(s string) (time.Duration, error) {
	if strings.HasSuffix(s, "d") {
		date := strings.TrimSuffix(s, "d")
		dateInt, err := strconv.Atoi(date)
		if err != nil {
			return 0, err
		}
		s = fmt.Sprintf("%dh", dateInt*24)
	}
	return time.ParseDuration(s)
}

func infof(msg string, v ...interface{}) {
	const prefix = "[INFO] "
	logf(prefix, msg, v...)
}

func errorf(msg string, v ...interface{}) {
	const prefix = "[ERROR] "
	logf(prefix, msg, v...)
}

func logf(prefix, msg string, v ...interface{}) {
	fmt.Printf(prefix + fmt.Sprintf(msg, v...) + "\n")
}
