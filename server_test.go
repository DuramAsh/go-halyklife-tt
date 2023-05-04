package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleProxy(t *testing.T) {
	router := gin.New()
	router.Any("/proxy", handleProxy)

	requestBody := RequestBody{
		Method: "GET",
		URL:    "https://google.com",
		Headers: map[string]interface{}{
			"Authentication": "Basic 832rdashzf812349",
			"Cookie":         "hsid=293hdfsalfad",
			"User-Agent":     "Mozilla/5.0 (Win64; x64)",
		},
	}
	requestBytes, _ := json.Marshal(requestBody)

	request := httptest.NewRequest("POST", "/proxy", bytes.NewBuffer(requestBytes))

	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, request)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	responseBody := &ResponseBody{}
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &responseBody)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, responseBody.StatusCode)
	fmt.Println(responseBody)
}
