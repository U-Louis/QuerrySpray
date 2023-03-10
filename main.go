package main

import (
    "bytes"
    "io/ioutil"
    "net/http"
    "strings"
    "errors"
    "log"
    "strconv"
    "context"
    "time"
    "fmt"

    "github.com/gin-gonic/gin"
)

type Sprayable struct {
    Method   string   `json:"method"`
    Uri      string   `json:"uri"`
    Protocol string   `json:"protocol"`
    Headers  []string `json:"headers"`
    Body     string   `json:"body"`
}

func main() {
    r := gin.Default()

    r.POST("/spray", func(c *gin.Context) {
        var sprayable Sprayable
    
        if err := c.BindJSON(&sprayable); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
    
        // Get the number of times to spray the request from the query parameter
        multipleStr := c.Query("multiple")
        if multipleStr == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Missing query parameter: multiple"})
            return
        }
    
        multiple, err := parseMultiple(multipleStr)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
    
        // Create the HTTP request based on the sprayable object
        req, err := http.NewRequest(sprayable.Method, sprayable.Uri, bytes.NewBuffer([]byte(sprayable.Body)))
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
    
        // Add headers to the HTTP request
        for _, header := range sprayable.Headers {
            parts := strings.Split(header, ":")
            if len(parts) != 2 {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid header: " + header})
                return
            }
            req.Header.Add(parts[0], strings.TrimSpace(parts[1]))
        }
    
        // Perform the HTTP request multiple times
        resp, err := performRequestMultipleTimes(http.DefaultClient, req, multiple)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
    
        if resp == nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "All requests failed"})
            return
        }
    
        // Return the first response received
        defer resp.Body.Close()

// Return the first response received
body, err := ioutil.ReadAll(resp.Body)
if err != nil {
    // log error and request
    log.Printf("Error: %s | Request: %v\n", err.Error(), req)
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
}
c.Data(http.StatusOK, "application/json", body)


    })   

    r.Run(":8085")
}

func performRequestMultipleTimes(client *http.Client, req *http.Request, multiple int) (*http.Response, error) {
    log.Printf("Request: %s %s, headers=%v, body=%s", req.Method, req.URL.String(), req.Header, req.Body)

    // Use a channel to make the requests concurrently and wait for the first response
    responseChan := make(chan *http.Response, multiple)
    errorChan := make(chan error, 1)

    // Use a context with a 60-second timeout
    ctx, cancel := context.WithTimeout(req.Context(), 60*time.Second)
    defer cancel()

    for i := 0; i < multiple; i++ {
        go func() {
            resp, err := performRequest(client, req.WithContext(ctx))
            if err != nil {
                // If an error occurs, send it on the error channel
                errorChan <- err
                return
            }

            // Send the response on the response channel
            responseChan <- resp
        }()
    }

    // Wait for the first response or error
    select {
    case resp := <-responseChan:
        log.Printf("Response: status=%d, headers=%v, body=%s", resp.StatusCode, resp.Header, resp.Body)
        return resp, nil
    case err := <-errorChan:
        return nil, err
    case <-ctx.Done():
        return nil, fmt.Errorf("timeout waiting for response after 60 seconds")
    }
}

func parseMultiple(multipleStr string) (int, error) {
    multiple, err := strconv.Atoi(multipleStr)
    if err != nil {
        return 0, errors.New("Invalid query parameter: multiple")
    }
    return multiple, nil
}

func performRequest(client *http.Client, req *http.Request) (*http.Response, error) {
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    log.Println(string(body))
    return resp, nil
}
