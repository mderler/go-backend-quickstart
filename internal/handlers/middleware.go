package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type contextKey string

const (
	userIDKey contextKey = "userID"
	todoIDKey contextKey = "todoID"
)

func userCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "id")
		if userID == "" {
			writeInvalidUserIdError(w, userID)
			return
		}
		id, err := strconv.ParseInt(userID, 10, 32)
		if err != nil {
			writeInvalidUserIdError(w, userID)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, int32(id))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func todoCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		todoID := chi.URLParam(r, "id")
		if todoID == "" {
			writeInvalidTodoIdError(w, todoID)
			return
		}
		id, err := strconv.ParseInt(todoID, 10, 32)
		if err != nil {
			writeInvalidTodoIdError(w, todoID)
			return
		}

		ctx := context.WithValue(r.Context(), todoIDKey, int32(id))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
