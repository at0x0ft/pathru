.PHONY: build
build:
	docker compose run --rm -u "$$(id -u):$$(id -g)" build github.com/at0x0ft/pathru

.PHONY: stat
stat:
	docker compose run --rm -u "$$(id -u):$$(id -g)" go version -m ./bin/*

.PHONY: update
update:
	docker compose run --rm -u "$$(id -u):$$(id -g)" go mod tidy

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
	docker compose run --rm -u "$$(id -u):$$(id -g)" go test -v ./...

.PHONY: lint_light
lint_light:
	docker compose run --rm -u "$$(id -u):$$(id -g)" go vet ./...

.PHONY: format
format:
	docker compose run --rm -u "$$(id -u):$$(id -g)" go fmt ./...
