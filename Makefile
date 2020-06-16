.PHONY: generate
generate:
	go mod vendor && \
	bazel run \
		--platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 \
		//:gazelle

.PHONY: test
test: generate
	bazel test //...

.PHONY: build
build: test
	bazel build \
		--platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 \
		//...

.PHONY: push
push: build
	bazel run :push_mail2wordpress

.PHONY: wordpress
wordpress: build
	go run .

.PHONY: spotify
spotify: build
	go run . spotify

.PHONY: local
local:
	go run . local-playlist
