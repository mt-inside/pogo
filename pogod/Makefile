ROOT := ..
HERE := pogod
PROD_IMAGE := pogod
include ../Makefile

.PHONY: pogod pogod-run

# Local
pogod:
	CGO_ENABLED=0
	go generate ./...
	go install ./...
pogod-run:
	go run main.go
