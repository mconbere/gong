package gong

import (
	"context"
	"net/http"
)

type ErrorHandler interface {
	ServeError(w http.ResponseWriter, r *http.Request, message string, status int) error
}

var notFoundHandlerContextKey = &ContextKey{name: "not found handler context key"}
var serverErrorHandlerContextKey = &ContextKey{name: "server error handler context key"}

func WithNotFoundHandler(ctx context.Context, h http.Handler) context.Context {
	return context.WithValue(ctx, notFoundHandlerContextKey, h)
}

func WithServerErrorHandler(ctx context.Context, h ErrorHandler) context.Context {
	return context.WithValue(ctx, serverErrorHandlerContextKey, h)
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	v, ok := ctx.Value(notFoundHandlerContextKey).(http.Handler)
	if !ok {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	v.ServeHTTP(w, r)
}

func ServerError(w http.ResponseWriter, r *http.Request, status int, message string) {
	ctx := r.Context()
	v, ok := ctx.Value(serverErrorHandlerContextKey).(ErrorHandler)
	if ok {
		w.WriteHeader(status)
		if err := v.ServeError(w, r, message, status); err == nil {
			return
		}
	}
	http.Error(w, message, status)
	return
}
