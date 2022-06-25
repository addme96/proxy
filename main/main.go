package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	apiURL, ok := os.LookupEnv("API_URL")
	if !ok {
		log.Fatal("required API_URL env variable not set")
	}
	proxy := NewProxy(apiURL, &http.Client{Timeout: 30 * time.Second})
	log.Println("Listening for requests at http://localhost:8080. Proxying to: " + apiURL)
	log.Fatal(http.ListenAndServe(":8080", Heartbeat("/healthz")(http.HandlerFunc(proxy.Proxy))))
}
