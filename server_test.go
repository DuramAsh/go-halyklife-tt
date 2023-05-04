package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetRequest(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(ProxyEndpoint))
	defer testServer.Close()

	reqBody := RequestBody{ // client -> proxy
		Method: "GET",
		URL:    "https://google.com",
		Headers: map[string]interface{}{
			"Authentication": "Basic 832rdashzf812349",
			"Cookie":         "hsid=293hdfsalfad",
			"User-Agent":     "Mozilla/5.0 (Win64; x64)",
		},
	}

	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequest(reqBody.Method, testServer.URL, bytes.NewReader(jsonBody)) // proxy -> service
	if err != nil {
		t.Error(err)
		return
	}

	res, err := http.DefaultClient.Do(req) // proxy <- service
	if err != nil {
		t.Error(err)
		return
	}
	defer res.Body.Close()
}
