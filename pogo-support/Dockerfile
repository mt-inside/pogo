# Run build-image to build it
# Multi-part to extract go binary in to scratch
FROM golang-pogo-builder as build
WORKDIR /go/src/github.com/mt-inside/pogo
COPY . .
RUN go generate ./pogo/...
ENV CGO_ENABLED=0
RUN go install ./pogo/...

FROM scratch
COPY --from=build /go/bin/pogo /

CMD ["/pogo"]
