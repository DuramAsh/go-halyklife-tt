package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
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

func handleProxy(ctx *gin.Context) {
	// Client -> Proxy
	C2PRequestBody := &RequestBody{}
	_ = ctx.BindJSON(&C2PRequestBody)
	// Proxy -> Service
	P2SRequest, err := http.NewRequest(C2PRequestBody.Method, C2PRequestBody.URL, ctx.Request.Body)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	// Proxy <- Service
	client := http.Client{}
	S2PResponse, err := client.Do(P2SRequest)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	// Client <- Proxy
	P2CResponseBody := &ResponseBody{
		ID:         777,
		StatusCode: S2PResponse.StatusCode,
		Headers:    convertHeaders(S2PResponse.Header),
		Length:     S2PResponse.ContentLength,
	}
	ctx.JSON(P2CResponseBody.StatusCode, P2CResponseBody)
	dumpData := make(map[string]interface{})
	dumpData["request"] = C2PRequestBody
	dumpData["response"] = P2CResponseBody
	Dump2JSON("dump.json", dumpData)
}

func convertHeaders(headers http.Header) map[string]interface{} {
	m := make(map[string]interface{})
	for k, v := range headers {
		m[k] = v
	}
	return m
}

func Dump2JSON(filename string, data map[string]interface{}) {
	file, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	jsonData := make([]map[string]interface{}, 0)
	decoder := json.NewDecoder(file)
	_ = decoder.Decode(&jsonData)
	file.Close()
	os.Remove(filename)

	file, _ = os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	jsonData = append(jsonData, data)
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	_ = encoder.Encode(jsonData)
	file.Close()
}

func main() {
	router := gin.Default()
	router.Any("/proxy", handleProxy)

	_ = router.Run()
}
