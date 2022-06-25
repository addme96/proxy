package main

import (
	"io"
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

	client := http.Client{Timeout: 30 * time.Second}
	log.Println("Listening for requests at http://localhost:8080. Proxying to: " + apiURL)
	log.Fatal(http.ListenAndServe(":8080", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path == "/healthz" {
			writer.WriteHeader(http.StatusOK)
			log.Printf("GET /healthz: Status 200 OK.")
			return
		}

		log.Printf(`Request:
URL: %s
Method: %s

`, request.URL.Path, request.Method)
		req, err := http.NewRequest(request.Method, apiURL+request.URL.Path, request.Body)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		req.Header = request.Header.Clone()
		resp, err := client.Do(req)
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(resp.StatusCode)
		_, _ = writer.Write(body)
		log.Printf(`Response:
StatusCode: %d
Body: %s

`, resp.StatusCode, string(body))
	})))
}
