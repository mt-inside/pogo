package(default_visibility = ["//visibility:public"])

load("@org_pubref_rules_protobuf//go:rules.bzl", "go_proto_library")

filegroup(
    name = "protos",
    srcs = [
        "pogo.proto",
    ],
)

go_proto_library(
    name = "go",
    protos = [":protos"],
    with_grpc = True,
    verbose = 1,
)
