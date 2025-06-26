package redis

import (
	"time"
)

type Config struct {
	Hosts        string
	Password     string
	Database     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}
