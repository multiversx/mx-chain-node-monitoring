package noderating

import "errors"

// ErrNilHTTPClient signals that a nil http client have been provided
var ErrNilHTTPClient = errors.New("nil http client")

// ErrEmptyPubKeys signals that no public keys have been provided in config
var ErrEmptyPubKeys = errors.New("no public keys provided in config")

// ErrEmptyApiUrl signals that an empty api url has been provided
var ErrEmptyApiUrl = errors.New("empty api url has been provided")
