BINARY = prom-exporter
ACT_URL = https://raw.githubusercontent.com/nektos/act/master/install.sh

ACT := $(shell command -v act)

.PHONY: test act test_sec test_ci docker_build goreleaser_build

test: test_sec test_ci

act:
	@[ ! -x "$(ACT)" ] && (curl --proto '=https' --tlsv1.2 -sSf $(ACT_URL) | sudo bash && sudo install ./bin/act /usr/local/bin/) || true

test_ci: act
	@act push --rm -j "test"

test_sec: act
	@act push --rm -j "security"

docker_build:
	@docker build --no-cache -t $(BINARY) .

goreleaser_build:
	@goreleaser build --rm-dir --single-target --clean --auto-snapshot

