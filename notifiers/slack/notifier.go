package slack

import (
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/node-monitoring/config"
	"github.com/ElrondNetwork/node-monitoring/data"
	"github.com/ElrondNetwork/node-monitoring/notifiers"
)

var log = logger.GetOrCreate("slackNotifier")

type payload struct {
	Text string `json:"text"`
}

// ArgsSlackNotifier defines the arguments needed to create a new slack notifier
type ArgsSlackNotifier struct {
	Config     *config.Slack
	HTTPClient notifiers.HTTPClient
}

type slackNotifier struct {
	url        string
	httpClient notifiers.HTTPClient
}

// NewSlackNotifier will create a new email notifier instance
func NewSlackNotifier(args *ArgsSlackNotifier) (*slackNotifier, error) {
	err := checkArgs(args)
	if err != nil {
		return nil, err
	}

	return &slackNotifier{
		url:        args.Config.URL,
		httpClient: args.HTTPClient,
	}, nil
}

func checkArgs(args *ArgsSlackNotifier) error {
	if args.Config.URL == "" {
		return ErrInvalidSlackURL
	}
	if args.HTTPClient == nil {
		return ErrNilHTTPClient
	}

	return nil
}

// PushMessage will push the notification
func (sn *slackNotifier) PushMessage(msg data.NotificationMessage) error {
	msgPayload := payload{
		Text: msg.Message,
	}

	return sn.httpClient.CallPostRestEndPoint(sn.url, "", msgPayload)
}

// GetID will return the identifier for slack notifier
func (sn *slackNotifier) GetID() string {
	return "Slack"
}

// IsInterfaceNil returns true if there is no value under the interface
func (sn *slackNotifier) IsInterfaceNil() bool {
	return sn == nil
}
