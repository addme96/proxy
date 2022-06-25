package main

import (
	"io"
	"log"
	"net/http"
)

type Proxy struct {
	TargetURL string
	client    *http.Client
}

func NewProxy(targetURL string, client *http.Client) *Proxy {
	return &Proxy{TargetURL: targetURL, client: client}
}

func (p Proxy) Proxy(writer http.ResponseWriter, request *http.Request) {
	resp, err := p.callAPI(request, p.TargetURL)
	if err != nil {
		log.Printf("callAPI error: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err = p.writeResponse(writer, resp); err != nil {
		log.Printf("writeResponse error: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (p Proxy) callAPI(request *http.Request, apiURL string) (*http.Response, error) {
	log.Printf("Request:\nURL: %s\nMethod: %s\n", request.URL.Path, request.Method)
	req, err := http.NewRequest(request.Method, apiURL+request.URL.Path, request.Body)
	if err != nil {
		return nil, err
	}
	req.Header = request.Header.Clone()
	return p.client.Do(req)
}

func (p Proxy) writeResponse(writer http.ResponseWriter, resp *http.Response) error {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	p.writeHeaders(writer, resp.Header)
	writer.WriteHeader(resp.StatusCode)
	_, err = writer.Write(body)
	if err != nil {
		return err
	}
	log.Printf("Response:\nStatusCode: %d\nBody: %s\n", resp.StatusCode, string(body))
	return nil
}

func (p Proxy) writeHeaders(writer http.ResponseWriter, header http.Header) {
	for key, vals := range header {
		for _, val := range vals {
			writer.Header().Add(key, val)
		}
	}
}
