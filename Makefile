.PHONY: help
help: # show list of all commands
	@grep -E '^[a-zA-Z_-]+:.*?# .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?# "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: test
test: TEST_DIR?=...
test: # run unit tests
	go test ./$(TEST_DIR) -count=1

.PHONY: lint
lint: # run linter
	golangci-lint run

.PHONY: fmt
fmt: # run formatter
	golangci-lint fmt

.PHONY: precommit
precommit: # run precommit checks
	$(MAKE) lint
	$(MAKE) test

.PHONY: todo
todo: # show list of all todos left in code
	@rg 'TODO' --glob '**/*.go' || echo 'All done!'
