# must ensure your go version >= 1.16
.PHONY: install
install:
	go install golang.org/x/tools/cmd/goimports@latest

.PHONY: tidy
tidy:
	@go mod tidy
	@$(foreach dir,$(shell go list -f {{.Dir}} ./...),goimports -w $(dir);)

.PHONY: test
test: tidy
	@go test ./...