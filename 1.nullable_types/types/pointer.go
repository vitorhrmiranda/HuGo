package types

import (
	"encoding/json"
	"net/http"
)

type PayloadWithPointerTypes struct {
	Count *int `json:"count"`
}

func NewAppThatAcceptsPointerTypes() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var payload PayloadWithPointerTypes
		if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
			http.Error(writer, "", http.StatusBadRequest)
			return
		}
		if payload.Count == nil {
			http.Error(writer, "", http.StatusUnprocessableEntity)
			return
		}
		*payload.Count++
		_ = json.NewEncoder(writer).Encode(payload)
	}
}
