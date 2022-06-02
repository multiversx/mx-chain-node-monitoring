package slack

import "errors"

// ErrInvalidSlackURL signals that an empty url has been provided
var ErrInvalidSlackURL = errors.New("empty slack url has been provided")

// ErrNilHTTPClient signals that a nil http client has been provided
var ErrNilHTTPClient = errors.New("nil http client")
