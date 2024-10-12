
test:
	@cd examples/functions && \
		CSV=true go run ../../cmd/nestcsv -a csv

	@cd examples/downstream && \
		go run ../../cmd/nestcsv
.PHONY: test
