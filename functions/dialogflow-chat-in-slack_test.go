package functions

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var handler http.Handler
var server *httptest.Server

func init() {
	handler = http.HandlerFunc(SimplestBotFunction)
	server = httptest.NewServer(handler)
}

func TestSimplestBotFunctionWithUrlVerification(t *testing.T) {

	t.Run("asserts url verified", func(t *testing.T) {
		body := map[string]interface{}{
			"token":     "Jhj5dZrVaK7ZwHHjRyZWjbDl",
			"challenge": "3eZbrw1aBm2rZgRNFdxV2595E9CY3gmdALWMmHkvFXO7tYXAYM8P",
			"type":      "url_verification",
		}
		bodyBytes, err := json.Marshal(body)
		request, err := http.NewRequest("POST", "/", bytes.NewReader(bodyBytes))
		response := httptest.NewRecorder()

		SimplestBotFunction(response, request)

		assert.Nil(t, err)

		assert.Nil(t, err)
		assert.Equal(t, response.Code, 200)
		assert.Equal(t, response.Body.String(), body["challenge"])

	})

}
