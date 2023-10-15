package routers

import (
	"expvar"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/bersennaidoo/mserve/application/services/transport/rest/boost-sales-api/handlers/debug/checkgrp"
	v1TestGrp "github.com/bersennaidoo/mserve/application/services/transport/rest/boost-sales-api/handlers/v1/testgrp"
	"github.com/bersennaidoo/mserve/business/web/v1/middleware"
	"github.com/bersennaidoo/mserve/library/web"
	"go.uber.org/zap"
)

func DebugStandardLibraryMux() *http.ServeMux {
	mux := http.NewServeMux()

	// Register all the standard library debug endpoints.
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())

	return mux
}

func DebugMux(build string, log *zap.SugaredLogger) http.Handler {

	mux := DebugStandardLibraryMux()

	cgh := checkgrp.Handlers{
		Build: build,
		Log:   log,
	}

	mux.HandleFunc("/debug/readiness", cgh.Readiness)
	mux.HandleFunc("/debug/liveness", cgh.Liveness)

	return mux
}

type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
}

func APIMux(cfg APIMuxConfig) *web.App {

	app := web.NewApp(
		cfg.Shutdown,
		middleware.Logger(cfg.Log),
		middleware.Errors(cfg.Log),
		middleware.Metrics(),
		middleware.Panics(),
	)

	v1(app, cfg)

	return app
}

func v1(app *web.App, cfg APIMuxConfig) {

	const version = "v1"

	tgh := v1TestGrp.Handlers{
		Log: cfg.Log,
	}

	app.Handle(http.MethodGet, version, "/test", tgh.Test)
}
