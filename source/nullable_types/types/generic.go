package types

import (
	"encoding/json"
	"net/http"

	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Float | constraints.Integer
}

type NullableNumber[T Number] struct {
	value *T
}

func (t NullableNumber[T]) Value() T {
	if t.value == nil {
		return *new(T)
	}
	return *t.value
}

func (t *NullableNumber[T]) Set(v T) {
	t.value = &v
}

func (t NullableNumber[T]) IsNull() bool {
	return t.value == nil
}

func (t *NullableNumber[T]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &t.value)
}

func (t NullableNumber[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.value)
}

type PayloadWithGenericNumberType struct {
	Count NullableNumber[float64] `json:"count"`
}

func NewAppThatAcceptsGenericNumberType() http.HandlerFunc {
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
