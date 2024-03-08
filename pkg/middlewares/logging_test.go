package middlewares

import (
	"log/slog"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

type stubHandler struct {
	counter int
}

func (h *stubHandler) handle(_ *fasthttp.RequestCtx) {
	h.counter++
}

func TestLoggingMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		ignorePaths    []string
		path           string
		wantLogMessage assert.BoolAssertionFunc
	}{
		{
			name:           "no ignore paths",
			ignorePaths:    nil,
			path:           "/metrics",
			wantLogMessage: assert.True,
		},
		{
			name:           "empty ignore paths",
			ignorePaths:    []string{},
			path:           "/metrics",
			wantLogMessage: assert.True,
		},
		{
			name:           "path in ignore paths",
			ignorePaths:    []string{"/metrics"},
			path:           "/metrics",
			wantLogMessage: assert.False,
		},
		{
			name:           "no path in ignore paths",
			ignorePaths:    []string{"/metrics"},
			path:           "/health",
			wantLogMessage: assert.True,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := new(strings.Builder)
			logger := slog.New(slog.NewTextHandler(sb, nil))
			mw := LoggingMiddleware(logger, tt.ignorePaths)
			next := stubHandler{}
			handler := mw(next.handle)
			ctx := new(fasthttp.RequestCtx)
			ctx.Request.SetRequestURI(tt.path)

			handler(ctx)

			assert.Equal(t, 1, next.counter)
			tt.wantLogMessage(t, len(sb.String()) != 0)
		})
	}
}
