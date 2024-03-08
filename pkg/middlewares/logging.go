package middlewares

import (
	"log/slog"
	"time"

	"github.com/valyala/fasthttp"
)

type Middleware func(next func(ctx *fasthttp.RequestCtx)) func(ctx *fasthttp.RequestCtx)

func LoggingMiddleware(logger *slog.Logger, ignorePaths []string) Middleware {
	ignorePathsMap := make(map[string]struct{}, len(ignorePaths))
	for _, path := range ignorePaths {
		ignorePathsMap[path] = struct{}{}
	}

	return func(next func(ctx *fasthttp.RequestCtx)) func(ctx *fasthttp.RequestCtx) {
		return func(ctx *fasthttp.RequestCtx) {
			if _, ok := ignorePathsMap[string(ctx.Path())]; !ok {
				beginTime := time.Now()
				defer func() {
					logger.Info("SERVER",
						"method", string(ctx.Method()),
						"uri", ctx.URI().String(),
						"status", ctx.Response.StatusCode(),
						"responseTime", time.Since(beginTime).String(),
					)
				}()
			}

			next(ctx)
		}
	}
}
