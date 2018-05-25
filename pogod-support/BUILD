package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@org_pubref_rules_protobuf//go:rules.bzl", "GRPC_COMPILE_DEPS")

go_binary(
    name = "pogod",
    srcs = [
        "main.go",
    ],
    deps = [
        "//proto:go",
    ] + GRPC_COMPILE_DEPS,
)
