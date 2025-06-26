package pkg

import (
	"context"
	"net/http"
	"net/http/pprof"
	"net/textproto"
	"runtime/debug"
	"sync/atomic"

	"github.com/LittleLe6owski/link-shortener/api"
	"github.com/LittleLe6owski/link-shortener/pkg/server"
	"github.com/go-openapi/runtime/middleware"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcProvider(app *server.App) *server.GRPCServer {
	options := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(unaryServerInterceptor(app.Logger())),
	}

	config := &server.GRPCServerConfig{
		Host:            app.Config().GrpcConfig.Host,
		Port:            app.Config().GrpcConfig.Port,
		ShutdownTimeout: app.Config().GrpcConfig.ShutdownTimeout,
		Options:         options,
	}

	return server.NewGRPCServer(config)
}

func GRPCGatewayHandlerProvider(app *server.App) *runtime.ServeMux {
	return runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(headerMatcher))
}

const (
	requestTagHeader = "X-Request-Tag"
)

func headerMatcher(key string) (string, bool) {
	key = textproto.CanonicalMIMEHeaderKey(key)
	switch key {
	case requestTagHeader:
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}

func unaryServerInterceptor(logger zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.With().Fields([]any{string(debug.Stack())}).Err(err)
				err = status.Error(codes.Internal, "internal error")
			}
		}()

		resp, err = handler(ctx, req)

		return resp, err
	}
}

func HttpServerProvider(app *server.App) *server.HTTPServer {
	config := &server.HTTPServerConfig{
		Host:            app.Config().HTTPConfig.Host,
		Port:            app.Config().HTTPConfig.Port,
		ShutdownTimeout: app.Config().GrpcConfig.ShutdownTimeout,
	}

	srv := server.NewHTTPServer(config)
	srv.Router().HandleFunc("/liveness", LivenessHandler())
	srv.Router().HandleFunc("/readiness", ReadinessHandler(app.ReadyState()))

	if app.Config().Service.Profiling {
		srv.Router().HandleFunc("/debug/pprof/", pprof.Index)
		srv.Router().HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		srv.Router().HandleFunc("/debug/pprof/profile", pprof.Profile)
		srv.Router().HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		srv.Router().HandleFunc("/debug/pprof/trace", pprof.Trace)
	}

	return srv
}

func AddSwagger(router *http.ServeMux) {
	router.Handle("/swagger.json", middleware.Spec("", []byte(api.SwaggerJSON), nil))
	router.Handle("/docs", middleware.SwaggerUI(middleware.SwaggerUIOpts{
		SpecURL: "/swagger.json",
	}, nil))
}

func LivenessHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func ReadinessHandler(state *atomic.Bool) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if state.Load() {
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
	}
}
