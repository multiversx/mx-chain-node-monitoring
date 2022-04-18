package email

import "errors"

// ErrInvalidEmailCredentials signals that invalid email credentials have been provided
var ErrInvalidEmailCredentials = errors.New("invalid email credentials")

// ErrEmptyEmailToList signals that no email to has been provided
var ErrEmptyEmailToList = errors.New("empty email to list")

// ErrInvalidEmailHostPort signals that an invalid email host port has been provided
var ErrInvalidEmailHostPort = errors.New("invalid email host port")
