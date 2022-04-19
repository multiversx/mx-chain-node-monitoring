package noderating

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/node-monitoring/clients"
	"github.com/ElrondNetwork/node-monitoring/common"
	"github.com/ElrondNetwork/node-monitoring/config"
	"github.com/ElrondNetwork/node-monitoring/data"
)

var log = logger.GetOrCreate("clients/nodeRating")

const (
	maxRating = 100
)

const (
	baseURL = "https://api.elrond.com"

	nodesBLSKeyPath = "/nodes/%s"
)

// ArgsNodeRating defines the arguments needed to create a new client
type ArgsNodeRating struct {
	Client clients.HTTPClient
	Config *config.NodeRating
}

type nodeRating struct {
	httpClient clients.HTTPClient
	lastValues map[string]float64
	firstRun   bool
	config     *config.NodeRating
}

// NewNodeRatingClient creates an instance of httpClient which is a wrapper for http.Client
func NewNodeRatingClient(args ArgsNodeRating) (*nodeRating, error) {
	err := checkArgs(args)
	if err != nil {
		return nil, err
	}

	return &nodeRating{
		httpClient: args.Client,
		lastValues: make(map[string]float64),
		firstRun:   true,
		config:     args.Config,
	}, nil
}

func checkArgs(args ArgsNodeRating) error {
	if check.IfNil(args.Client) {
		return ErrNilHTTPClient
	}
	if len(args.Config.PubKeys) == 0 {
		return fmt.Errorf("no public keys provided in config")
	}
	if args.Config.Threshold <= 0 {
		return fmt.Errorf("%w: invalid node rating threshold, provided %.2f", common.ErrInvalidValue, args.Config.Threshold)
	}

	return nil
}

func (hcw *nodeRating) GetEvent() (data.NotificationMessage, error) {
	if hcw.firstRun == true {
		log.Info("First run. Will not trigger any event.")
		hcw.firstRun = false

		return hcw.handleFirstRun()
	}

	return hcw.handleEvents()
}

func (hcw *nodeRating) handleEvents() (data.NotificationMessage, error) {
	event := data.NotificationMessage{Level: common.InfoEvent}

	nodes, err := hcw.fetchAPINodesByBLSKey()
	if err != nil {
		return event, err
	}

	msg := ""
	for _, node := range nodes {
		// TODO: evaluate simply if the rating drops under a specified threshold
		diff := node.TempRating - hcw.lastValues[node.Bls]
		if diff >= 0 {
			continue
		}

		changePercentage := math.Abs(diff) / float64(maxRating)
		changePercentage = changePercentage * 100
		if changePercentage > hcw.config.Threshold {
			event.Level = common.CriticalEvent
			nodeMsg := fmt.Sprintf(
				"NodeName: %s - TempRating decreased with %.1f percent, current value: %.2f, last value: %.2f\n",
				node.Name,
				hcw.config.Threshold,
				node.TempRating,
				hcw.lastValues[node.Bls],
			)
			msg = msg + nodeMsg
		}

		hcw.lastValues[node.Bls] = node.TempRating
	}

	event.Message = msg

	return event, nil
}

func (hcw *nodeRating) handleFirstRun() (data.NotificationMessage, error) {
	nodes, err := hcw.fetchAPINodesByBLSKey()
	if err != nil {
		return data.NotificationMessage{}, err
	}

	for _, node := range nodes {
		hcw.lastValues[node.Bls] = node.TempRating
	}

	return data.NotificationMessage{Level: common.NoEvent}, nil
}

func (hcw *nodeRating) fetchAPINodesByBLSKey() ([]clients.APINode, error) {
	nodes := make([]clients.APINode, 0)

	for _, pubKey := range hcw.config.PubKeys {
		path := fmt.Sprintf(nodesBLSKeyPath, pubKey)
		responseBodyBytes, err := hcw.httpClient.CallGetRestEndPoint(baseURL, path)
		if err != nil {
			return nil, err
		}

		var response clients.APINode
		err = json.Unmarshal(responseBodyBytes, &response)
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, response)
	}

	return nodes, nil
}

func (hcw *nodeRating) GetID() string {
	return "NodeRating"
}

// IsInterfaceNil returns true if there is no value under the interface
func (hcw *nodeRating) IsInterfaceNil() bool {
	return hcw == nil
}
