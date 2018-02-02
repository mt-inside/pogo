BUILD_IMAGE := bazel
DEV_IMAGE := bazel-dev
PROD_IMAGE := pogod

# Local
build:
	bazel build //main:pogod

run:
	bazel run //main:pogod

# Production docker image
image:
	docker build -t $(PROD_IMAGE) .

image-run:
	docker run -p8080:8080 $(PROD_IMAGE)

# "Docker monad"
.builder-image:
	docker build -t $(BUILD_IMAGE) build-image

.dev-image: .builder-image
	docker build -t $(DEV_IMAGE) -f build-image/Dockerfile.dev build-image

docker-%: .builder-image
	docker run \
	    -v $(shell pwd):/app \
	    $(BUILD_IMAGE) \
	    make $*

docker-dev: .dev-image
	docker run -ti -v $(shell pwd):/app $(DEV_IMAGE)
