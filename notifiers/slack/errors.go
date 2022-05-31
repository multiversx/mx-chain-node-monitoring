package slack

import "errors"

// ErrInvalidSlackURL signals that an empty url has been provided
var ErrInvalidSlackURL = errors.New("empty slack url has been provided")
