package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// RequestFields is for holding request details
type RequestFields struct {
	Body    string                 `json:"body"`
	Cookies []*http.Cookie         `json:"cookies"`
	Path    string                 `json:"path"`
	Params  map[string][]string    `json:"params"`
	Method  string                 `json:"method"`
	Headers map[string]interface{} `json:"headers"`
}

// HTTPFields is for holding protocol details
type HTTPFields struct {
	Version string               `json:"version"`
	TLS     *tls.ConnectionState `json:"tls"`
}

// MetaFields is for holding connection information
type MetaFields struct {
	Host       string `json:"host"`
	RemoteAddr string `json:"remote_addr"`
	Timestamp  string `json:"timestamp"`
	UUID       string `json:"uuid"`
}

// Reflection holds the reflection information
type Reflection struct {
	HTTP    HTTPFields    `json:"http"`
	Request RequestFields `json:"data"`
	Meta    MetaFields    `json:"meta"`
}

// InitReflection is used for initialising the uninitialisable stuff
// in the Reflection struct
func InitReflection() *Reflection {
	reflection := Reflection{}
	reflection.Request.Headers = make(map[string]interface{}, 0)
	return &reflection
}

// Handler is a http.HandleFunc to pass into a HTTP server
func Handler(w http.ResponseWriter, r *http.Request) {
	reflection := InitReflection()

	reflection.HTTP.TLS = r.TLS
	reflection.HTTP.Version = r.Proto

	if body, err := ioutil.ReadAll(r.Body); err != nil {
		panic(err)
	} else {
		reflection.Request.Body = string(body)
	}
	reflection.Request.Cookies = r.Cookies()
	for key, value := range r.Header {
		reflection.Request.Headers[key] = value
	}
	reflection.Request.Method = r.Method
	reflection.Request.Params = r.URL.Query()
	reflection.Request.Path = r.URL.Path

	reflection.Meta.UUID = w.Header().Get("Request-ID")
	reflection.Meta.Host = r.Host
	reflection.Meta.RemoteAddr = r.RemoteAddr
	reflection.Meta.Timestamp = time.Now().Format("2006-01-02T15:04:05-0700")

	w.Header().Set("Content-Type", "application/json")

	if response, err := json.Marshal(reflection); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		w.Write([]byte("\"something went wrong\""))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%s", response)))
	}
}
