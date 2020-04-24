package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func sendToSlack(data SlackRequest) error {
	if data.Channel == "" {
		data.Channel = GetEnvSlackChannel()
	}
	if data.UserName == "" {
		data.UserName = GetEnvSlackUser()
	}

	slackBody, _ := json.Marshal(data)
	req, err := http.NewRequest(http.MethodPost, GetEnvSlackURL(), bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	cli := &http.Client{Timeout: 10 * time.Second}
	resp, err := cli.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return err
	}

	body := buf.String()
	if body != "ok" {
		return fmt.Errorf("unknown response from Slack API: [%s]", body)
	}
	return nil
}

type SlackRequest struct {
	Channel     string            `json:"channel,omitempty"`
	UserName    string            `json:"username,omitempty"`
	Attachments []SlackAttachment `json:"attachments"`
}

type SlackAttachment struct {
	Color      string   `json:"color,omitempty"`
	Title      string   `json:"title,omitempty"`
	TitleLink  string   `json:"title_link,omitempty"`
	MarkdownIn []string `json:"mrkdwn_in,omitempty"`
	Text       string   `json:"text,omitempty"`
}
