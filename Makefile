POGO_REL := github.com/mt-inside/pogo
CONTAINER_POGO_ABS := /go/src/$(POGO_REL)
BUILD_IMAGE := golang-pogo-builder
DEV_IMAGE := golang-pogo-builder-dev

.PHONY: image image-run .builder-image .dev-image docker-dev

image:
	docker build -t $(PROD_IMAGE) -f ./Dockerfile $(ROOT)

image-run: image
	docker run --network pogo_net $(PROD_IMAGE)
#TODO: make a dockercompose you idiot. Separate dockerfiles in subdirs.
#	Write separate make then work out how to factor it out later


# "Docker monad"
.builder-image:
	docker build -t $(BUILD_IMAGE) $(ROOT)/build-image

.dev-image: .builder-image
	docker build -t $(DEV_IMAGE) -f $(ROOT)/build-image/Dockerfile.dev $(ROOT)/build-image

docker-%: .builder-image
	docker run \
	    -v $(shell realpath $(ROOT)):$(CONTAINER_POGO_ABS) \
	    $(BUILD_IMAGE) \
	    /bin/sh -c "cd $(CONTAINER_POGO_ABS)/$(HERE) && make $*"

docker-dev: .dev-image
	docker run -ti -v $(shell realpath $(ROOT)):$(CONTAINER_POGO_ABS) $(DEV_IMAGE)


pogo:
	CGO_ENABLED=0
	vgo generate github.com/mt-inside/pogo/cmd/pogo
	vgo install github.com/mt-inside/pogo/cmd/pogo

pogod:
	CGO_ENABLED=0
	vgo generate github.com/mt-inside/pogo/cmd/pogod
	vgo install github.com/mt-inside/pogo/cmd/pogod

pogod-image:
	docker build -t hack --build-arg binary=pogod -f ./Dockerfile .
