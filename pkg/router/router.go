package router

import (
	"github.com/dimfeld/httptreemux/v5"
	"net/http"
)

func New() *httptreemux.ContextMux {

	newRouter := httptreemux.NewContextMux()

	newRouter.GET("/health/status", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return newRouter
}
