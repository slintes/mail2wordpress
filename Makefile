.PHONY: build
build:
	go mod vendor && \
	bazel run \
		--platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 \
		//:gazelle && \
	bazel build \
		--platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 \
		//...

.PHONY: push
push: build
	bazel run :push_mail2wordpress