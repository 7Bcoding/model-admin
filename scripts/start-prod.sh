#!/bin/bash
export $(cat .env.production | grep -v '^#' | xargs)
go run main.go 