package multiplexer

import (
	"log/slog"
	"net/http"
)

type Middleware struct {
	logger *slog.Logger
}

func NewMiddleware(logger *slog.Logger) *Middleware {
	return &Middleware{
		logger: logger,
	}
}

func (mw *Middleware) printReq(r *http.Request) {
	mw.logger.Info(
		"request",
		slog.Any("remote address", r.RemoteAddr),
		slog.Any("user-agent", r.UserAgent()),
	)
}

func (mw *Middleware) CliLogger(f http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		mw.printReq(r)
		f.ServeHTTP(rw, r)
	}
}
