package notify

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type SlackSettings struct {
	WebhookUrl string `json:"webhook_url"`
}

type SlackNotifier struct {
	Settings *SlackSettings
}

// There is a whole lot of loosely documented things you can to with a Slack
// message. This could be expanded later on to include all the posiblities.
type SlackPayload struct {
	Text        string                   `json:"text"`
	Attachments []map[string]interface{} `json:"attachments,omitempty"`
}

func (s *SlackNotifier) Notify(text string) error {
	body, err := json.Marshal(SlackPayload{
		Text: "GOSSM check failed",
		Attachments: []map[string]interface{}{
			map[string]interface{}{
				"title": "Server information:",
				"text":  text,
				"color": "danger",
			},
		},
	})
	if err != nil {
		return err
	}

	resp, err := http.Post(s.Settings.WebhookUrl, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New(string(data))
	}
	return nil
}

func (ss *SlackSettings) Validate() error {
	if ss.WebhookUrl == "" {
		return errors.New("You must provide a webhook for your slack message. See https://api.slack.com/incoming-webhooks for more information")
	}
	webHook, err := url.Parse(ss.WebhookUrl)
	if err != nil {
		return err
	}
	if webHook.Hostname() != "hooks.slack.com" {
		return errors.New("You must provide a webhook for your slack message. See https://api.slack.com/incoming-webhooks for more information")
	}
	return err
}

func (s *SlackNotifier) String() string {
	return fmt.Sprintf("Slack: with webhookURL %s", s.Settings.WebhookUrl)
}
