package app

import (
	"context"
	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
}

func (a *App) Initialize(ctx context.Context) {
	a.Router = mux.NewRouter()
	a.setRouters()
}