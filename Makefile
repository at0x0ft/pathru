repository_root := "$$(git rev-parse --show-toplevel)"

docker_compose_run := docker compose run --rm -u "$$(id -u):$$(id -g)"

common_build_command := go build -o ./bin -ldflags "-s -w" ./...
common_build_options := CGO_ENABLED=0
common_test_command := go test -v ./...
common_format_command := go fmt ./...

.PHONY: setup_pre_commit
setup_pre_commit:
	cd $(repository_root) && \
	ln -svf ../../.githooks/pre-commit .git/hooks/pre-commit

.PHONY: build
build:
	cd $(repository_root) && \
	$(docker_compose_run) -e '$(common_build_options)' $(common_build_command)

.PHONY: stat
stat:
	cd $(repository_root) && \
	$(docker_compose_run) go version -m ./bin/*

.PHONY: pkg_resolve
pkg_resolve:
	cd $(repository_root) && \
	$(docker_compose_run) go mod tidy

.PHONY: clean
clean:
	cd $(repository_root) && \
	docker compose down -v && \
	rm -rf ./bin/*

.PHONY: cache_clear
cache_clear:
	cd $(repository_root) && \
	sudo rm -rf ./.go_build/* && \
	git checkout HEAD -- ./.go_build

.PHONY: test
test:
	cd $(repository_root) && \
	$(docker_compose_run) $(common_test_command)

.PHONY: lint
lint:
	cd $(repository_root) && \
	$(docker_compose_run) go vet ./...

.PHONY: format
format:
	cd $(repository_root) && \
	$(docker_compose_run) $(common_format_command)

# === command(s) for CI ===

.PHONY: ci_build
ci_build:
	cd $(repository_root) && \
	$(common_build_options) $(common_build_command)

.PHONY: ci_test
ci_test:
	cd $(repository_root) && \
	$(common_test_command)

# === command(s) for git pre-commit hooks ===

.PHONY: git_pre_commit_hook_format
git_pre_commit_hook_format:
	cd $(repository_root) && \
	$(docker_compose_run) -T $(common_format_command)
