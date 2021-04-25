package server

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"time"
)

func analyticMiddleware(isTestMode bool, excludedURI []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		skipCurrentURI := false
		for _, uri := range excludedURI {
			if ctx.Request.RequestURI == uri {
				skipCurrentURI = true
				break
			}
		}

		var bodyBytes []byte
		if ctx.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(ctx.Request.Body)
			ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		ctx.Next()

		latency := float64(0)
		end := time.Now()
		latency = float64(end.Sub(start).Nanoseconds()/1e4) / 100.0

		if !skipCurrentURI && ctx.Writer.Status() != 0 && !isTestMode {
			log.Info().
				Str("client_ip", ctx.ClientIP()).
				Int("status", ctx.Writer.Status()).
				Str("method", ctx.Request.Method).
				Str("endpoint", ctx.Request.URL.Path).
				Str("query", ctx.Request.URL.RawQuery).
				Str("body", string(bodyBytes)).
				Str("time", start.Format(time.RFC3339)).
				Float64("latency", latency).
				Str("user_agent", ctx.Request.UserAgent()).
				Msg("API REQUEST")
		}
	}
}
