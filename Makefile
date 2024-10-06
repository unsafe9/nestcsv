
test:
	@cd examples/functions && \
		rm -rf ./go ./json ./ue5 && \
		CSV=true go run ../../cmd/nestcsv -a csv

.PHONY: test
