
test:
	@cd examples && \
		rm -r ./go ./json ./ue5 && \
		go run ../cmd/nestcsv

.PHONY: test
