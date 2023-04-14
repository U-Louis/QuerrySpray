#!/bin/bash

## Remove the old binary
if [ -f "main" ]; then
  rm main
  fi

## Build the new binary
#go build main.go

## Build for linux
GOOS=linux GOARCH=amd64 go build main.go

## Run the new binary
./main