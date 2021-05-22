build-all:
	goreleaser --snapshot --rm-dist

release:
	goreleaser release --rm-dist
