
test:
	@cd examples/functions && \
		CSV=true go run ../../cmd/nestcsv -a csv

	@cd examples/downstream && \
		go run ../../cmd/nestcsv
.PHONY: test

build-local:
	goreleaser build --snapshot --clean

release:
	#export GITHUB_TOKEN=...
	git tag -a v$(VERSION) -m "Release v$(VERSION)"
	git push origin v$(VERSION)
	goreleaser release --clean
.PHONY: release
