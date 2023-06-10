.PHONY: build
build:
	docker compose run --rm -u "$$(id -u):$$(id -g)" build github.com/at0x0ft/pathru/cmd/pathru

.PHONY: stat
stat:
	docker compose run --rm -u "$$(id -u):$$(id -g)" go version -m ./bin/*

.PHONY: clean
clean:
	docker compose down -v && \
	rm -rf ./bin/*

.PHONY: cache_clear
cache_clear:
	sudo rm -rf ./.go_build/* && \
	git checkout HEAD -- ./.go_build
