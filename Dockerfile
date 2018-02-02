# Run build-image to build it
# Multi-part to extract go binary in to scratch
FROM golang-pogo-builder as build
WORKDIR /go/src/github.com/mt-inside/pogo
COPY . .
RUN go generate ./pogod
ENV CGO_ENABLED=0
RUN go install ./pogod

FROM scratch
COPY --from=build /go/bin/pogod /

EXPOSE 50001
CMD ["/pogod"]
