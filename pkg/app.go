package pkg

import (
	"context"
	"os"
	"sync"

	"github.com/LittleLe6owski/link-shortener/internal/config"
	"github.com/rs/zerolog"
)

type App struct {
	servers                   []Server
	logger                    zerolog.Logger
	config                    *config.Config
	preHooks                  []Hook
	postHooks                 []Hook
	afterServersShutdownHooks []Hook
}

type Hook func(ctx context.Context, app *App) error

func NewApp(config config.Config) (*App, error) {
	return &App{
		config: &config,
		logger: zerolog.New(os.Stdout).With().Timestamp().Logger(),
	}, nil
}

// AddServer adds server to app. App servers will be run during App.Run execution.
func (a *App) AddServer(srv Server) {
	a.servers = append(a.servers, srv)
}

func (a *App) Config() *config.Config {
	return a.config
}

func (a *App) Logger() zerolog.Logger {
	return a.logger
}

func (a *App) AddPreHook(hook Hook) {
	a.preHooks = append(a.preHooks, hook)
}

func (a *App) AddPostHook(hook Hook) {
	a.postHooks = append(a.postHooks, hook)
}

func (a *App) AddAfterServersShutdownHook(hook Hook) {
	a.afterServersShutdownHooks = append(a.afterServersShutdownHooks, hook)
}

func (a *App) Run(ctx context.Context) error { //nolint:cyclop // by design many if
	ctxWithCancel, cancel := context.WithCancel(ctx)
	defer cancel()
	sysCtx := newSystemContext(ctxWithCancel, a.config.Service.GracefulShutdownDelay)

	// pre hooks
	for _, hook := range a.preHooks {
		if err := hook(sysCtx, a); err != nil {
			a.logger.Err(ErrPreHook)
			return err
		}
	}

	// start servers
	for _, srv := range a.servers {
		if err := srv.Start(); err != nil {
			a.logger.Err(ErrStartServer)
			return err
		}
	}

	// post hooks
	var postErr error
	for _, hook := range a.postHooks {
		if postErr = hook(sysCtx, a); postErr != nil {
			a.logger.Err(ErrPostHook)
			cancel()
			break
		}
	}
	if postErr == nil {
		a.logger.Info().Msg("app started successfully")
	}

	// listen for context cancel
	<-sysCtx.Done()

	// stop servers in parallel
	wg := new(sync.WaitGroup)
	for _, srv := range a.servers {
		wg.Add(1)
		go func(srv Server) {
			if err := srv.Stop(); err != nil {
				a.logger.Err(ErrStopServer)
			}
			wg.Done()
		}(srv)
	}
	wg.Wait()

	for _, hook := range a.afterServersShutdownHooks {
		if err := hook(context.Background(), a); err != nil {
			a.logger.Err(ErrShutDownHook)
		}
	}

	a.logger.Info().Msg("app stopped successfully")

	return postErr
}
