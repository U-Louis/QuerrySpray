package main

import (
    "bytes"
    "io/ioutil"
    "net/http"
    "strings"
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

        // Make the HTTP request
        client := &http.Client{}
        res, err := client.Do(req)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        defer res.Body.Close()

        // Read the response body
        body, err := ioutil.ReadAll(res.Body)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        // Return the response
        c.JSON(http.StatusOK, gin.H{
            "statusCode": res.StatusCode,
            "headers":    res.Header,
            "body":       string(body),
        })
    })

    r.Run(":8085")
}
