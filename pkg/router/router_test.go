package router

import (
	"github.com/appleboy/gofight/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestNew(t *testing.T) {

	app := New()
	r := gofight.New()

	t.Run("health test", func(t *testing.T) {
		r.GET("/health/status").
			SetDebug(true).
			Run(app, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusOK, r.Code)
			})
	})
}
