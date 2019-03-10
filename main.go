// creates the version file

//go:generate go run ./generators/version/main.go

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"sync"

	"github.com/google/uuid"
)

const serverTLSKeyFile = "server.key"
const serverTLSCertFile = "server.crt"

func main() {
	var port int
	var version bool
	var waiter sync.WaitGroup

	httpDone := make(chan error)
	httpsDone := make(chan error)
	flag.IntVar(&port, "port", 8080, "defines the port that the reflection server listens on - the https listener adds 1 to this for the port")
	flag.BoolVar(&version, "version", false, "display the version")
	flag.Parse()

	if version {
		fmt.Printf("%s-%s", Version, Commit)
		return
	}

	defer log.Println("stopping the reflection server")
	log.Println("starting the reflection server...")
	http.HandleFunc("/", withLogging(Handler))

	// handle the http
	waiter.Add(1)
	go listenHTTP(port, &httpDone)
	// handle the https
	useTLS := shouldEnableTLS()
	if useTLS {
		waiter.Add(1)
		go listenHTTPS(port+1, &httpsDone)
	}
	// wait for both to terminate
	go handleShutdown(&httpsDone, &httpDone, &waiter)

	waiter.Wait()
}

func getCurrentDirectory() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return cwd
}

func handleShutdown(httpsDone *chan error, httpDone *chan error, waiter *sync.WaitGroup) {
	for {
		select {
		case httpsErr := <-*httpsDone:
			if httpsErr != nil {
				log.Printf("https server terminated with: %v\n", httpsErr)
			}
			waiter.Done()
		case httpErr := <-*httpDone:
			if httpErr != nil {
				log.Printf("http server terminated with: %v\n", httpErr)
			}
			waiter.Done()
		}
	}
}

func listenHTTP(port int, done *chan error) {
	listenAddress := fmt.Sprintf("0.0.0.0:%v", port)
	log.Printf("attempting to listen for http on %s...\n", listenAddress)
	*done <- http.ListenAndServe(listenAddress, nil)
}

func listenHTTPS(port int, done *chan error) {
	cwd := getCurrentDirectory()
	listenAddress := fmt.Sprintf("0.0.0.0:%v", port)
	log.Printf("attempting to listen for https on %s...\n", listenAddress)
	*done <- http.ListenAndServeTLS(listenAddress, path.Join(cwd, "/server.crt"), path.Join(cwd, "server.key"), nil)
}

func shouldEnableTLS() bool {
	if keyInfo, err := os.Lstat(path.Join(getCurrentDirectory(), serverTLSKeyFile)); os.IsNotExist(err) {
		log.Printf("%s not found - generate with `make tlscerts` for https to be available\n", serverTLSKeyFile)
	} else if certInfo, err := os.Lstat(path.Join(getCurrentDirectory(), "server.crt")); os.IsNotExist(err) {
		log.Printf("%s not found - generate with `make tlscerts` for https to be available\n", serverTLSCertFile)
	} else if !keyInfo.IsDir() && !certInfo.IsDir() {
		log.Printf("%s and %s found - tls will be enabled\n", serverTLSKeyFile, serverTLSCertFile)
		return true
	} else {
		log.Printf("%s and %s found but they were not quite as expected - generate with `make tlscerts` for https to be available", serverTLSKeyFile, serverTLSCertFile)
	}
	return false
}

func withLogging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		protocolStub := "http"
		if r.TLS != nil {
			protocolStub = fmt.Sprintf("%ss", protocolStub)
		}
		log.Printf("%s %s [id:%s] %s://%s%s\n", r.Method, r.Proto, requestID, protocolStub, r.Host, r.RequestURI)
		w.Header().Add("Request-ID", requestID)
		next.ServeHTTP(w, r)
	}
}
