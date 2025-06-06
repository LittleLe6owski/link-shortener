package pkg

import "errors"

var (
	ErrPreHook      = errors.New("error occurred executing app pre hook")
	ErrStartServer  = errors.New("failed to start app server")
	ErrPostHook     = errors.New("error occurred executing app post hook")
	ErrStopServer   = errors.New("failed to stop app server")
	ErrShutDownHook = errors.New("failed after servers shutdown hook")
)
