FROM golang:1.10-stretch as builder

RUN apt update && apt install -y unzip

# vgo for better package management
RUN go get -u golang.org/x/vgo

# protoc compiler and stdlib types
RUN wget -O /tmp/protoc.zip https://github.com/google/protobuf/releases/download/v3.5.1/protoc-3.5.1-linux-x86_64.zip
RUN unzip -d /usr/local /tmp/protoc.zip
# golang code generation backend for protoc
RUN go get -u github.com/golang/protobuf/protoc-gen-go


FROM builder as build
ARG binary

WORKDIR /go/src/github.com/mt-inside/pogo
COPY . .

ENV CGO_ENABLED=0
RUN vgo generate github.com/mt-inside/pogo/cmd/$binary
RUN vgo install github.com/mt-inside/pogo/cmd/$binary


FROM scratch

COPY --from=build /go/bin/$binary /

ENTRYPOINT ["/pogod"]
