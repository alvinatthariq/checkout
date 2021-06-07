.PHONY: run-tests
run-tests:
	@go test -v -tags dynamic `go list ./...` -cover

.PHONY: build
build:
	@go build main.go

.PHONY: run
run: run-tests build
	@./main
