ROOT := ..
HERE := pogo
PROD_IMAGE := pogo
include ../Makefile

.PHONY: pogo pogod-run

# Local
pogo:
	CGO_ENABLED=0
	go generate ./...
	go install ./...
pogo-run:
	go run main.go
