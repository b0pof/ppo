package response

import (
	"encoding/json"
	"net/http"
)

func OK(w http.ResponseWriter, data interface{}) {
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(data)
}

func Unauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode("Авторизация отсутствует")
}

type errorResponse struct {
	Msg string `json:"error"`
}

func BadRequest(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(w).Encode(errorResponse{Msg: msg})
}

func Internal(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(errorResponse{Msg: err.Error()})
}

func Forbidden(w http.ResponseWriter, msg *string) {
	w.WriteHeader(http.StatusForbidden)
	if msg != nil {
		_ = json.NewEncoder(w).Encode(errorResponse{Msg: *msg})
	}
}
