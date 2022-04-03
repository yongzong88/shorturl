package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const HOST = "http://localhost:9000"

// 測試加短網址，透過 http
func Test_AddComplete(t *testing.T) {

	ct := time.Now().Add(time.Second * 100)
	reqbody := map[string]string{
		"url":      "https://www.hinet.net",
		"expireAt": ct.Format(time.RFC3339),
	}
	jsonStr, _ := json.Marshal(reqbody)
	req, _ := http.NewRequest("POST", HOST+"/api/v1/urls", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, error := client.Do(req)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	t.Log("HOST=", req.Host)
	t.Log("response Status:", response.Status)
	t.Log("response Headers:", response.Header)
	body, _ := ioutil.ReadAll(response.Body)
	t.Log("response Body:", string(body))

	assert.Equal(t, http.StatusOK, response.StatusCode)
}

// 測試加短網址，透過 handle 直接操作
func Test_AddHandle(t *testing.T) {
	router := router()

	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	ct := time.Now().Add(time.Second * 10)
	reqbody := map[string]string{
		"url":      "https://www.hinet.net",
		"expireAt": ct.Format(time.RFC3339),
	}
	jsonStr, _ := json.Marshal(reqbody)
	req, _ := http.NewRequest("POST", "/api/v1/urls", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)
	t.Log(w.Body.String())

	assert.Equal(t, http.StatusOK, w.Code)
}

// 測試加短網址，然後拜訪短網址
func Test_AddThenVisit(t *testing.T) {
	router := router()

	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件

	ct := time.Now().Add(time.Second * 100)
	t.Log("CT", ct)
	testUrl := "https://www.google.com"
	reqbody := map[string]string{
		"url":      testUrl,
		"expireAt": ct.Format(time.RFC3339),
	}
	jsonStr, _ := json.Marshal(reqbody)

	req, _ := http.NewRequest("POST", "/api/v1/urls", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 取得回應
	var response map[string]string = make(map[string]string)
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, err, nil)

	t.Log("Short URL =", response["id"])

	req, _ = http.NewRequest("GET", HOST+"/"+response["id"], nil)
	client := &http.Client{}
	result, err := client.Do(req)
	assert.Equal(t, err, nil)
	defer result.Body.Close()

	assert.Equal(t, result.StatusCode, 200)
	t.Log("response Body:", result.Request.URL)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Equal(t, result.Request.URL.String(), testUrl)
}

// 測試加短網址，過時之後再拜訪短網址
func Test_AddThenVisitAfterTimeout(t *testing.T) {
	router := router()

	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件

	// 有效時間 3 秒
	ct := time.Now().Add(3 * time.Second)
	t.Log("CT", ct)
	testUrl := "https://www.google.com"
	reqbody := map[string]string{
		"url":      testUrl,
		"expireAt": ct.Format(time.RFC3339),
	}
	jsonStr, _ := json.Marshal(reqbody)

	req, _ := http.NewRequest("POST", "/api/v1/urls", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 取得回應
	var response map[string]string = make(map[string]string)
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, err, nil)

	t.Log("Short URL =", response["id"])

	// 	暫停 5 秒
	time.Sleep(5 * time.Second)

	req, _ = http.NewRequest("GET", HOST+"/"+response["id"], nil)
	client := &http.Client{}
	result, err := client.Do(req)
	assert.Equal(t, err, nil)
	defer result.Body.Close()

	assert.Equal(t, result.StatusCode, 404)
}

// 測試不存在的短網址
func Test_VisitNonExist(t *testing.T) {
	req, _ := http.NewRequest("GET", HOST+"/ZZ", nil)
	client := &http.Client{}
	result, err := client.Do(req)
	assert.Equal(t, err, nil)
	defer result.Body.Close()

	assert.Equal(t, result.StatusCode, 404)
}

// 測試加短網址(Handle)的效能
func Benchmark_AddHandle(b *testing.B) {

	router := router()
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	ct := time.Now().Add(time.Second * 10)
	reqbody := map[string]string{
		"url":      "https://www.hinet.net",
		"expireAt": ct.Format(time.RFC3339),
	}
	jsonStr, _ := json.Marshal(reqbody)

	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("POST", "/api/v1/urls", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
	}
}

// 測試加短網址(Handle)，然後拜訪短網址的效能，可能會被視為攻擊行為
// func Benchmark_AddThenMultiVisit(b *testing.B) {
// 	router := router()

// 	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件

// 	ct := time.Now().Add(time.Second * 100)
// 	testUrl := "http://www.hinet.net"
// 	reqbody := map[string]string{
// 		"url":      testUrl,
// 		"expireAt": ct.Format(time.RFC3339),
// 	}
// 	jsonStr, _ := json.Marshal(reqbody)

// 	req, _ := http.NewRequest("POST", "/api/v1/urls", bytes.NewBuffer(jsonStr))
// 	req.Header.Set("Content-Type", "application/json")

// 	router.ServeHTTP(w, req)

// 	// 取得回應
// 	var response map[string]string = make(map[string]string)
// 	json.Unmarshal(w.Body.Bytes(), &response)

// 	for i := 0; i < b.N; i++ {
// 		req, _ = http.NewRequest("GET", HOST+"/"+response["id"], nil)
// 		client := &http.Client{}
// 		result, _ := client.Do(req)
// 		defer result.Body.Close()
// 	}
// }
