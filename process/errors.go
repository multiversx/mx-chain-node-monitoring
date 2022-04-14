package process

import "errors"

// ErrNilClient signals that a nil client have been provided
var ErrNilClient = errors.New("nil client")

// ErrNilPusher signals that a nil pusher instance have been provided
var ErrNilPusher = errors.New("nil pusher instance")
