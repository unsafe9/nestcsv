
test:
	@cd examples && go run ../cmd/nestcsv
test2:
	@cd examples/test2 && go run ../../cmd/nestcsv

.PHONY: test
