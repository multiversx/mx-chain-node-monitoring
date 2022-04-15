package noderating

import "errors"

// ErrNilHTTPClient signals that a nil http client have been provided
var ErrNilHTTPClient = errors.New("nil http client")

// ErrEmptyPubKeys signals that no public keys have been provided in config
var ErrEmptyPubKeys = errors.New("no public keys provided in config")
