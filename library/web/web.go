package web

import (
	"context"
	"net/http"
	"os"

	"github.com/dimfeld/httptreemux/v5"
)

type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

type App struct {
	mux      *httptreemux.ContextMux
	shutdown chan os.Signal
}

func NewApp(shutdown chan os.Signal) *App {

	mux := httptreemux.NewContextMux()

	return &App{
		mux:      mux,
		shutdown: shutdown,
	}
}
