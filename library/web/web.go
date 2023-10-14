package web

import (
	"context"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/dimfeld/httptreemux/v5"
	"go.opentelemetry.io/otel/trace"
)

type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

type App struct {
	mux      *httptreemux.ContextMux
	shutdown chan os.Signal
	mw       []Middleware
}

func NewApp(shutdown chan os.Signal, mw ...Middleware) *App {

	mux := httptreemux.NewContextMux()

	return &App{
		mux:      mux,
		shutdown: shutdown,
		mw:       mw,
	}
}

func (a *App) SignalShutdown() {

	a.shutdown <- syscall.SIGTERM
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}

func (a *App) Handle(method string, group string, path string, handler Handler, mw ...Middleware) {

	handler = wrapMiddleware(mw, handler)

	handler = wrapMiddleware(a.mw, handler)

	h := func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		span := trace.SpanFromContext(ctx)

		v := Values{
			TraceID: span.SpanContext().TraceID().String(),
			Now:     time.Now(),
		}
		ctx = context.WithValue(ctx, key, &v)

		if err := handler(ctx, w, r); err != nil {
			a.SignalShutdown()
			return
		}
	}

	finalPath := path
	if group != "" {
		finalPath = "/" + group + path
	}
	a.mux.Handle(method, finalPath, h)
}
