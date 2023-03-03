package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Sprayable struct {
	URL          string            `json:"url"`
	Method       string            `json:"method"`
	Params       map[string]string `json:"params"`
	Authorization string           `json:"authorization"`
	Headers      map[string]string `json:"headers"`
	Body         string            `json:"body"`
}

type SprayRequest struct {
	Sprayable    Sprayable `json:"sprayable"`
	Multiplicator int       `json:"multiplicator"`
}

type SprayResponse struct {
	Result string `json:"result"`
}

func (s *Sprayable) buildRequest() (*http.Request, error) {
	values := url.Values{}
	for k, v := range s.Params {
		values.Set(k, v)
	}

	req, err := http.NewRequest(s.Method, s.URL, strings.NewReader(s.Body))
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = values.Encode()

	if s.Authorization != "" {
		req.Header.Set("Authorization", s.Authorization)
	}

	for k, v := range s.Headers {
		req.Header.Set(k, v)
	}

	return req, nil
}

func (s *SprayRequest) spray(resultChan chan<- string) {
	for i := 0; i < s.Multiplicator; i++ {
		query, err := s.Sprayable.buildRequest()
		if err != nil {
			resultChan <- fmt.Sprintf("Failed to build request: %v", err)
			return
		}

		client := &http.Client{Timeout: 5 * time.Second}

		// Log the request being sent
		fmt.Printf("Sending request %d: %s %s\n", i+1, s.Sprayable.Method, s.Sprayable.URL)

		resp, err := client.Do(query)
		if err != nil {
			resultChan <- fmt.Sprintf("Failed to get response: %v", err)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			resultChan <- fmt.Sprintf("Failed to read response: %v", err)
			return
		}

		resultChan <- string(body)
	}
}

func sprayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var sr SprayRequest
	err := json.NewDecoder(r.Body).Decode(&sr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode request: %v", err), http.StatusBadRequest)
		return
	}

	resultChan := make(chan string)
	go sr.spray(resultChan)

	// Wait for the first response
	response := <-resultChan

	// Close the channel to prevent leaks
	close(resultChan)

	// Send the first response back to the client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SprayResponse{Result: response})
}

func main() {
	http.HandleFunc("/spray", sprayHandler)
	http.ListenAndServe(":8085", nil)
}
