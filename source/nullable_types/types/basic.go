package types

import (
	"encoding/json"
	"net/http"
)

type PayloadWithBasicTypes struct {
	Count int `json:"count"`
}

func NewAppThatAcceptsBasicTypes() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var payload PayloadWithBasicTypes
		if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
			http.Error(writer, "", http.StatusBadRequest)
			return
		}
		payload.Count++
		_ = json.NewEncoder(writer).Encode(payload)
	}
}
