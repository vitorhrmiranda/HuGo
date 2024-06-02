package types

import (
	"encoding/json"
	"net/http"
)

type NullableInt struct {
	value *int
}

func (t NullableInt) Value() int {
	if t.value == nil {
		return 0
	}
	return *t.value
}

func (t *NullableInt) Set(v int) {
	t.value = &v
}

func (t NullableInt) IsNull() bool {
	return t.value == nil
}

func (t *NullableInt) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &t.value)
}

func (t NullableInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.value)
}

type PayloadWithCustomTypes struct {
	Count NullableInt `json:"count"`
}

func NewAppThatAcceptsCustomTypes() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var payload PayloadWithCustomTypes
		if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
			http.Error(writer, "", http.StatusBadRequest)
			return
		}
		if payload.Count.IsNull() {
			http.Error(writer, "", http.StatusUnprocessableEntity)
			return
		}
		payload.Count.Set(payload.Count.Value() + 1)
		_ = json.NewEncoder(writer).Encode(payload)
	}
}
