package domain

import "errors"

var ErrLinkAlreadyExists = errors.New("link already exists")
var ErrLinkNotFound = errors.New("link not found")
