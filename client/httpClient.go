package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"sync"
	"time"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/node-monitoring/common"
	"github.com/ElrondNetwork/node-monitoring/config"
	"github.com/ElrondNetwork/node-monitoring/data"
)

// TODO: refactor to separate http client wrapper from event logic

var log = logger.GetOrCreate("client/httpClient")

const (
	firstRunValue    float64 = -1
	maxRating                = 10000000
	minReqTimeoutSec         = 1
)

const (
	baseURL = "https://api.elrond.com"

	nodesIdentifierPath = "/nodes?identity=%s"
	nodesBLSKeyPath     = "/nodes/%s"
)

type httpClientWrapper struct {
	httpClient    *http.Client
	mutHTTPClient sync.RWMutex
	lastValues    map[string]float64
	firstRun      bool
	config        *config.NodeRating
}

// HTTPClientWrapperArgs defines the arguments needed to create a new client
type HTTPClientWrapperArgs struct {
	ReqTimeoutSec int
	Config        *config.NodeRating
}

// NewHTTPClientWrapper creates an instance of httpClient which is a wrapper for http.Client
func NewHTTPClientWrapper(args HTTPClientWrapperArgs) (*httpClientWrapper, error) {
	if args.ReqTimeoutSec < minReqTimeoutSec {
		return nil, fmt.Errorf("%w, provided: %d, minimum: %d", common.ErrInvalidValue, args.ReqTimeoutSec, minReqTimeoutSec)
	}

	httpClient := http.DefaultClient
	httpClient.Timeout = time.Duration(args.ReqTimeoutSec) * time.Second

	return &httpClientWrapper{
		httpClient: httpClient,
		lastValues: make(map[string]float64),
		firstRun:   true,
		config:     args.Config,
	}, nil
}

func (hcw *httpClientWrapper) GetEvent() (data.NotificationMessage, error) {
	if hcw.firstRun == true {
		log.Info("First run. Will not trigger any event.")
		hcw.firstRun = false

		return hcw.handleFirstRun()
	}

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
		if changePercentage > hcw.config.Threshold {
			event.Level = common.CriticalEvent
			nodeMsg := fmt.Sprintf(
				"%s: TempRating decreased with %.1f percent, current value: %.2f, last value: %.2f\n",
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

func (hcw *httpClientWrapper) handleFirstRun() (data.NotificationMessage, error) {
	nodes, err := hcw.fetchAPINodesByBLSKey()
	if err != nil {
		return data.NotificationMessage{}, err
	}

	for _, node := range nodes {
		hcw.lastValues[node.Bls] = node.TempRating
	}

	return data.NotificationMessage{Level: common.NoEvent}, nil
}

func (hcw *httpClientWrapper) fetchAPINodesByBLSKey() ([]APINode, error) {
	nodes := make([]APINode, 0)

	for _, pubKey := range hcw.config.PubKeys {
		path := fmt.Sprintf(nodesBLSKeyPath, pubKey)
		responseBodyBytes, err := hcw.callGetRestEndPoint(baseURL, path)
		if err != nil {
			return nil, err
		}

		var response APINode
		err = json.Unmarshal(responseBodyBytes, &response)
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, response)
	}

	return nodes, nil
}

// callGetRestEndPoint calls an external end point
func (hcw *httpClientWrapper) callGetRestEndPoint(
	address string,
	path string,
) ([]byte, error) {
	req, err := http.NewRequest("GET", address+path, nil)
	if err != nil {
		return nil, err
	}

	userAgent := "Elrond Node Monitoring / 1.0.0 <Requesting data from api>"
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", userAgent)

	resp, err := hcw.httpClient.Do(req)
	if err != nil {
		if isTimeoutError(err) {
			return nil, err
		}

		return nil, err
	}

	defer func() {
		errNotCritical := resp.Body.Close()
		if errNotCritical != nil {
			log.Warn("base process GET: close body", "error", errNotCritical.Error())
		}
	}()

	return ioutil.ReadAll(resp.Body)
}

func isTimeoutError(err error) bool {
	if err, ok := err.(net.Error); ok && err.Timeout() {
		return true
	}

	return false
}
