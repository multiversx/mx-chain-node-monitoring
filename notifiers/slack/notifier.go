package slack

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/ElrondNetwork/node-monitoring/config"
	"github.com/ElrondNetwork/node-monitoring/data"
)

type payload struct {
	Text string `json:"text"`
}

// ArgsSlackNotifier defines the arguments needed to create a new slack notifier
type ArgsSlackNotifier struct {
	Config *config.Slack
}

type slackNotifier struct {
	config *config.Slack
}

// NewSlackNotifier will create a new email notifier instance
func NewSlackNotifier(args ArgsSlackNotifier) (*slackNotifier, error) {
	err := checkArgs(args)
	if err != nil {
		return nil, err
	}

	return &slackNotifier{
		config: args.Config,
	}, nil
}

func checkArgs(args ArgsSlackNotifier) error {
	if args.Config.URL == "" {
		return ErrInvalidSlackURL
	}

	return nil
}

// PushMessage will push the notification
func (sn *slackNotifier) PushMessage(msg data.NotificationMessage) error {
	return sn.push(msg)
}

func (sn *slackNotifier) push(msg data.NotificationMessage) error {
	data := payload{
		Text: msg.Message,
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", sn.config.URL, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// GetID will return the identifier for slack notifier
func (sn *slackNotifier) GetID() string {
	return "Slack"
}

// IsInterfaceNil returns true if there is no value under the interface
func (sn *slackNotifier) IsInterfaceNil() bool {
	return sn == nil
}
