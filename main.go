package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

type RequestBody struct {
	Method  string                 `json:"method"`
	URL     string                 `json:"url"`
	Headers map[string]interface{} `json:"headers"`
}

type ResponseBody struct {
	ID         uint                   `json:"id"`
	StatusCode int                    `json:"statusCode"`
	Headers    map[string]interface{} `json:"headers"`
	Length     int64                  `json:"length"`
}

func ProxyEndpoint(w http.ResponseWriter, r *http.Request) {
	reqFromClient := RequestBody{}
	_ = json.NewDecoder(r.Body).Decode(&reqFromClient)                                 // get json body of the client's request to req var
	reqToService, err := http.NewRequest(reqFromClient.Method, reqFromClient.URL, nil) // creating new request, based on the client's request
	if err != nil {
		log.Fatal("cannot create req from proxy to service", http.StatusBadRequest)
		return
	}

	for k, v := range reqFromClient.Headers { // setting all the headers from prev req to the new req
		reqToService.Header.Set(k, fmt.Sprintf("%v", v))
	}
	client := &http.Client{}
	serviceResponse, err := client.Do(reqToService)
	if err != nil {
		log.Fatal("failed to send request to service", http.StatusBadRequest)
		return
	}

	resp := &ResponseBody{
		ID:         uint(rand.Intn(100)),
		StatusCode: serviceResponse.StatusCode,
		Headers:    convertHeaders(serviceResponse.Header),
		Length:     serviceResponse.ContentLength,
	}
	json.NewEncoder(w).Encode(resp)
}

func convertHeaders(header http.Header) map[string]interface{} {
	headers := make(map[string]interface{})
	for k, v := range header {
		headers[k] = v
	}
	return headers
}

func main() {
	server := http.Server{}
	fmt.Println("Server is running...")
	http.HandleFunc("/proxy", ProxyEndpoint)
	_ = server.ListenAndServe()
}
