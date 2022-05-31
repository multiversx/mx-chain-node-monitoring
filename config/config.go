package config

// GeneralConfig will hold the configs
type GeneralConfig struct {
	Flags     *FlagsConfig
	General   *General
	Notifiers *Notifiers
	Alarms    *Alarms
}

// General holds the general configuration
type General struct {
	TriggerIntervalSec int
}

// Notifiers holds the configuration for notifiers
type Notifiers struct {
	Email *Email
	Slack *Slack
}

// Alarms holds the configuration for the alarms defined
type Alarms struct {
	NodeRating *NodeRating
}

// NodeRating holds the configuration for node rating alarm
type NodeRating struct {
	Threshold float64
	ApiUrl    string
	PubKeys   []string
}

// Email holds the configuration for email notifier
type Email struct {
	Enabled       bool
	EmailHost     string
	EmailPort     int
	EmailUsername string
	EmailPassword string
	From          string
	To            []string
}

// Slack holds the configuration for slack notifier
type Slack struct {
	Enabled bool
	URL     string
}

// FlagsConfig holds the values for CLI flags
type FlagsConfig struct {
	GeneralConfigPath string
}
