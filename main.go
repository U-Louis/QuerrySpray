package main

import (
    "bytes"
    "io"
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
    Method       string   `json:"method"`
    Uri          string   `json:"uri"`
    Protocol     string   `json:"protocol"`
    Headers      []string `json:"headers"`
    Body         string   `json:"body"`
    ResponseType string   `json:"responseType"`
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
    
        if resp.StatusCode >= 200 && resp.StatusCode < 300 {
            // Pass headers
            for key, values := range resp.Header {
                for _, value := range values {
                    c.Header(key, value)
                }
            }
    
            // Check if the response is chunked
if resp.TransferEncoding != nil && strings.Contains(strings.Join(resp.TransferEncoding, ","), "chunked") {
    // Set the Content-Type header based on the response type
    c.Header("Content-Type", sprayable.ResponseType)

    // Create a buffer to accumulate the chunks
    buf := bytes.NewBuffer(nil)

    // Copy the response body in chunks
    _, err := io.Copy(buf, resp.Body)
    if err != nil {
        log.Printf("Error reading response body: %s | Request: %v\n", err.Error(), req)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Log the response body
    log.Printf("Response Body: %s", buf.String())

    // Write the accumulated data from the buffer to the response
    _, err = c.Writer.Write(buf.Bytes())
    if err != nil {
        log.Printf("Error writing response: %s | Request: %v\n", err.Error(), req)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
} else {
    // Respond with the original response body and status code
    c.Status(resp.StatusCode)

    // Create a buffer to capture the response body
    buf := bytes.NewBuffer(nil)

    // Copy the response body to the buffer and log its contents
    _, err := io.Copy(buf, resp.Body)
    if err != nil {
        log.Printf("Error reading response body: %s | Request: %v\n", err.Error(), req)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    log.Printf("Response Body: %s", buf.String())

    // Write the accumulated data from the buffer to the response
    _, err = c.Writer.Write(buf.Bytes())
    if err != nil {
        log.Printf("Error writing response: %s | Request: %v\n", err.Error(), req)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
}

    
            defer resp.Body.Close()
        } else {
            c.JSON(resp.StatusCode, gin.H{"error": fmt.Sprintf("Received status code: %d", resp.StatusCode)})
            return
        }
    })

    r.Run(":8085")
}


func performRequestMultipleTimes(client *http.Client, req *http.Request, multiple int) (*http.Response, error) {
    // Use a channel to make the requests concurrently and wait for the first response
    responseChan := make(chan *http.Response, multiple)
    errorChan := make(chan error, 1)

    // Use a context with a 60-second timeout
    ctx, cancel := context.WithTimeout(req.Context(), 60*time.Second)
    defer cancel()

    for i := 0; i < multiple; i++ {
        // create a new request object for each request
        log.Printf("Request: %s %s, headers=%v, body=%s", req.Method, req.URL.String(), req.Header, req.Body)
        newReq, err := http.NewRequest(req.Method, req.URL.String(), req.Body)
        if err != nil {
            return nil, err
        }
        newReq.Header = req.Header.Clone()

        go func() {
            resp, err := performRequest(client, newReq.WithContext(ctx))
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

    return resp, nil
}

