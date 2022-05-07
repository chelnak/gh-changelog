tag:
	@git tag -a $(version) -m "Release $(version)"
	@git push --follow-tags

lint:
	@golangci-lint run ./...

build:
	@WORKINGDIR=$(pwd) goreleaser build --snapshot --rm-dist --single-target

.PHONY: mocks
mocks:
	@mockery --all
