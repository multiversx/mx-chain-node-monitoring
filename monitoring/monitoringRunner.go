package monitoring

import (
	"os"
	"os/signal"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/node-monitoring/clients"
	noderating "github.com/ElrondNetwork/node-monitoring/clients/nodeRating"
	"github.com/ElrondNetwork/node-monitoring/config"
	"github.com/ElrondNetwork/node-monitoring/notifiers/email"
	"github.com/ElrondNetwork/node-monitoring/notifiers/slack"
	"github.com/ElrondNetwork/node-monitoring/process"
)

var log = logger.GetOrCreate("monitoring")

const reqTimeoutSec = 10

type monitoringRunner struct {
	config *config.GeneralConfig
}

// NewMonitoringRunner create a new notifierRunner instance
func NewMonitoringRunner(cfgs *config.GeneralConfig) (*monitoringRunner, error) {
	if cfgs == nil {
		return nil, ErrNilConfigs
	}

	return &monitoringRunner{
		config: cfgs,
	}, nil
}

// Start will trigger the main flow
func (mr *monitoringRunner) Start() error {
	clientArgs := clients.HTTPClientWrapperArgs{
		ReqTimeoutSec: reqTimeoutSec,
	}
	httpClientWrapper, err := clients.NewHTTPClientWrapper(clientArgs)
	if err != nil {
		return err
	}

	nodeRatingArgs := noderating.ArgsNodeRating{
		Client: httpClientWrapper,
		Config: mr.config.Alarms.NodeRating,
	}
	nodeRatingClient, err := noderating.NewNodeRatingClient(nodeRatingArgs)
	if err != nil {
		return err
	}

	notifyProcessor := process.NewNotifyProcessor()

	if mr.config.Notifiers.Slack.Enabled {
		argsSlackNotifier := slack.ArgsSlackNotifier{Config: mr.config.Notifiers.Slack}
		slackNotifier, err := slack.NewSlackNotifier(argsSlackNotifier)
		if err != nil {
			return err
		}
		notifyProcessor.AddNotifier(slackNotifier)
	}

	if mr.config.Notifiers.Email.Enabled {
		argsEmailNotifier := email.ArgsEmailNotifier{Config: mr.config.Notifiers.Email}
		emailNotifier, err := email.NewEmailNotifier(argsEmailNotifier)
		if err != nil {
			return err
		}
		notifyProcessor.AddNotifier(emailNotifier)
	}

	argsEventsProcessor := process.ArgsEventsProcessor{
		Pusher:             notifyProcessor,
		TriggerInternalSec: mr.config.General.TriggerIntervalSec,
	}
	eventsProcessor, err := process.NewEventsProcessor(argsEventsProcessor)
	if err != nil {
		return err
	}
	eventsProcessor.AddClients(nodeRatingClient)

	eventsProcessor.Run()

	err = waitForGracefulShutdown(eventsProcessor)
	if err != nil {
		return err
	}

	return nil
}

func waitForGracefulShutdown(
	processor processorHandler,
) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit

	log.Info("closing components...")

	err := processor.Close()
	if err != nil {
		return err
	}

	return nil
}
