package testgrp

import (
	"context"
	"errors"
	"math/rand"
	"net/http"

	v1 "github.com/bersennaidoo/mserve/business/web/v1"
	"github.com/bersennaidoo/mserve/library/web"
	"go.uber.org/zap"
)

type Handlers struct {
	Log *zap.SugaredLogger
}

func (h Handlers) Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	if n := rand.Intn(100); n%2 == 0 {
		return v1.NewRequestError(errors.New("trusted error"), http.StatusBadRequest)
	}

	status := struct {
		Status string
	}{
		Status: "OK",
	}

	return web.Respond(ctx, w, status, http.StatusOK)
}
