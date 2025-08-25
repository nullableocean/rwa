package handlers

import (
	"errors"
	"net/http"
	"strings"
)

const (
	UserCtxKey = "user_id"

	TokenHeader = "Authorization"
	TokenPrefix = "Token "
)

func GetUserIdFromRequestCtx(r *http.Request) (int64, error) {
	ctx := r.Context()
	id, ok := ctx.Value(UserCtxKey).(int64)
	if !ok {
		return -1, errors.New("Not found user in request context")
	}

	return id, nil
}

func GetTokenFromRequest(r *http.Request) string {
	val := r.Header.Get(TokenHeader)
	if val == "" || !strings.HasPrefix(val, TokenPrefix) {
		return ""
	}

	return val[len(TokenPrefix):]
}

func badJsonError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Error json"))
}
