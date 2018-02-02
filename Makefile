POGO_SRC_REL := github.com/mt-inside/pogo
POGO_SRC_ABS := /go/src/$(POGO_SRC_REL)
POGOD_SRC_REL := $(POGO_SRC_REL)/pogod
POGOD_SRC_ABS := /go/src/$(POGOD_SRC_REL)
BUILD_IMAGE := golang-pogo-builder
DEV_IMAGE := golang-pogo-builder-dev
PROD_IMAGE := pogod

# Local
build:
	go generate ./pogod
	go install ./pogod

run:
	go run ./pogod/main.go

# Production docker image
image: .builder-image
	docker build -t $(PROD_IMAGE) .

image-run: image
	docker run -p8080:8080 $(PROD_IMAGE)

# "Docker monad"
.builder-image:
	docker build -t $(BUILD_IMAGE) build-image

.dev-image: .builder-image
	docker build -t $(DEV_IMAGE) -f build-image/Dockerfile.dev build-image

docker-%:
	docker run \
	    -v $(shell pwd):$(POGO_SRC_ABS) \
	    $(BUILD_IMAGE) \
	    /bin/sh -c "cd $(POGO_SRC_ABS) && make $*"

docker-dev: .dev-image
	docker run -ti -v $(shell pwd):$(POGO_SRC_ABS) $(DEV_IMAGE)
