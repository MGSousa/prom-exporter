BINARY = prom-exporter
ACT_URL = https://raw.githubusercontent.com/nektos/act/master/install.sh

ACT := $(shell command -v act)

test_ci:
	@[ ! -x "$(ACT)" ] && (curl --proto '=https' --tlsv1.2 -sSf $(ACT_URL) | sudo bash) || true
	@act push --rm

docker_build:
	@docker build --no-cache -t $(BINARY) .

goreleaser_build:
	@goreleaser build --rm-dir --single-target --clean --auto-snapshot

