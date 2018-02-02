workspace(name = "org_uk_empty_pogod")

# ================================================================

git_repository(
    name = "io_bazel_rules_go",
    remote = "https://github.com/bazelbuild/rules_go.git",
    tag = "0.9.0",
)

load("@io_bazel_rules_go//go:def.bzl", "go_repositories")

go_repositories()

# ================================================================

git_repository(
    name = "org_pubref_rules_protobuf",
    #tag = "v0.8.1",
    # need https://github.com/pubref/rules_protobuf/pull/159
    tag = "master",
    remote = "https://github.com/pubref/rules_protobuf.git",
)

load("@org_pubref_rules_protobuf//go:rules.bzl", "go_proto_repositories")

go_proto_repositories()
