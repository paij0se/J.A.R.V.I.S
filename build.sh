#!/bin/bash

GOOS=linux go build -o jarvis main.go
zip jarvis-linux.zip jarvis
