package link

import "errors"

var ErrInvalidRequest = errors.New("invalid request")
var ErrUnexpected = errors.New("unexpected error occurred")
var ErrMyItemAlreadyExists = errors.New("my-item already exists")
var ErrMyItemNotFound = errors.New("my-item not found")
