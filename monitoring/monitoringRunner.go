package monitoring

import (
	"os"
	"os/signal"

	"github.com/ElrondNetwork/node-monitoring/client"
	"github.com/ElrondNetwork/node-monitoring/config"
	"github.com/ElrondNetwork/node-monitoring/notifiers/email"
	"github.com/ElrondNetwork/node-monitoring/process"
)

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
	clientArgs := client.HTTPClientWrapperArgs{
		ReqTimeoutSec: reqTimeoutSec,
		Config:        mr.config.Alarms.NodeRating,
	}
	connector, err := client.NewHTTPClientWrapper(clientArgs)
	if err != nil {
		return err
	}

	notifyProcessor := process.NewNotifyProcessor()

	if mr.config.Notifiers.Email.Enabled {
		argsEmailNotifier := email.ArgsEmailNotifier{Config: mr.config.Notifiers.Email}
		emailNotifier, err := email.NewEmailNotifier(argsEmailNotifier)
		if err != nil {
			return err
		}
		notifyProcessor.AddNotifier(emailNotifier)
	}

	argsEventsProcessor := process.ArgsEventsProcessor{
		Client:             connector,
		Pusher:             notifyProcessor,
		TriggerInternalSec: mr.config.General.TriggerIntervalSec,
	}
	eventsProcessor, err := process.NewEventsProcessor(argsEventsProcessor)
	if err != nil {
		return err
	}

	eventsProcessor.Run()

	err = waitForGracefulShutdown(eventsProcessor)
	if err != nil {
		return err
	}

	return nil
}

func waitForGracefulShutdown(
	processor eventsHandler,
) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit

	err := processor.Close()
	if err != nil {
		return err
	}

	return nil
}
