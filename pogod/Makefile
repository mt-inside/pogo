ROOT := ..
HERE := pogod
PROD_IMAGE := pogod
include ../Makefile

.PHONY: pogod pogod-run

# Local
debug:
	CGO_ENABLED=0
	go generate ./...
	go install -gcflags "-N -l" ./...
prod:
	CGO_ENABLED=0
	go generate ./...
	go install -ldflags "-w" ./...
pogod-run:
	go run main.go
