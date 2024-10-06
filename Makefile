
test:
	@cd examples/functions && \
		rm -rf ./go ./json ./ue5 && \
		CSV=true go run ../../cmd/nestcsv -a csv

	@cd examples/downstream && \
		go run ../../cmd/nestcsv
.PHONY: test
