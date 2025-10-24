package request

import (
	"encoding/json"
	"net/http"
)

func ParseBody[T any](r *http.Request) (T, error) {
	var data T

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return data, err
	}

	return data, nil
}
