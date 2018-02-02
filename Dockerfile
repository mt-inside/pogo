# Run build-image to build it
# Multi-part to extract go binary in to scratch
FROM bazel as build
WORKDIR /app
COPY . .
RUN bazel build //main:pogod

FROM scratch
COPY --from=build /app/bazel-bin/main/linux_amd64_stripped/pogod /

EXPOSE 8080
CMD ["/pogod"]
