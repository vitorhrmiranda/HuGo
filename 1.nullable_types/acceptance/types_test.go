package acceptance_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"nullable/types"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAcceptance(t *testing.T) {
	t.Run("with basic type", func(t *testing.T) { Acceptance(t, types.NewAppThatAcceptsBasicTypes()) })
	t.Run("with pointer type", func(t *testing.T) { Acceptance(t, types.NewAppThatAcceptsPointerTypes()) })
	t.Run("with custom type", func(t *testing.T) { Acceptance(t, types.NewAppThatAcceptsCustomTypes()) })
}

func Acceptance(t *testing.T, app http.HandlerFunc) {
	t.Helper()

	testCases := []struct {
		wantStatusCode      int
		requestPayload      string
		wantResponsePayload string
	}{
		{http.StatusOK, "{\"count\":1}", "{\"count\":2}\n"},
		{http.StatusBadRequest, "", "\n"},
		{http.StatusUnprocessableEntity, "{}", "\n"},
		{http.StatusUnprocessableEntity, "{\"count\":null}", "\n"},
	}
	for _, testCase := range testCases {
		t.Run("", func(t *testing.T) {
			server := httptest.NewServer(app)
			t.Cleanup(func() { server.Close() })

			client, err := http.NewRequest(http.MethodPost, server.URL, strings.NewReader(testCase.requestPayload))
			require.NoError(t, err)
			response, err := http.DefaultClient.Do(client)
			require.NoError(t, err)
			t.Cleanup(func() { response.Body.Close() })

			var buf bytes.Buffer
			_, err = io.Copy(&buf, response.Body)
			require.NoError(t, err)

			assert.Equal(t, testCase.wantStatusCode, response.StatusCode)
			assert.Equal(t, testCase.wantResponsePayload, buf.String())
		})
	}
}
