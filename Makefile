docker_compose_run := docker compose run --rm -u "$$(id -u):$$(id -g)"

common_build_command := go build -o ./bin -ldflags "-s -w" ./...
common_build_options := CGO_ENABLED=0
common_test_command := go test -v ./...
common_format_command := go fmt ./...

.PHONY: build
build:
	$(docker_compose_run) -e '$(common_build_options)' $(common_build_command)

.PHONY: stat
stat:
	$(docker_compose_run) go version -m ./bin/*

.PHONY: pkg_update
pkg_update:
	$(docker_compose_run) go mod tidy

.PHONY: clean
clean:
	docker compose down -v && \
	rm -rf ./bin/*

.PHONY: cache_clear
cache_clear:
	sudo rm -rf ./.go_build/* && \
	git checkout HEAD -- ./.go_build

.PHONY: test
test:
	$(docker_compose_run) $(common_test_command)

.PHONY: lint
lint:
	$(docker_compose_run) go vet ./...

.PHONY: format
format:
	$(docker_compose_run) $(common_format_command)

# === commands for CI ===

.PHONY: ci_build
ci_build:
	$(common_build_options) $(common_build_command)

.PHONY: ci_test
ci_test:
	$(common_test_command)

# .PHONY: ci_lint
# ci_lint:
#	:	# ci_lint command defined in GHA workflow file.

.PHONY: ci_format
ci_format:
	$(common_format_command)
