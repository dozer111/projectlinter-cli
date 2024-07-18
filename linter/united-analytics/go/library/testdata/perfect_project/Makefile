UNITED_LTD_DIRECTORY:=${HOME}/.united_ltd
LINT_CONFIG_DIRECTORY:=${UNITED_LTD_DIRECTORY}/golangci-lint-config
LINT_CONFIG_FILE:=${UNITED_LTD_DIRECTORY}/golangci-lint-config/.golangci.yml

.PHONY: check-golang
check-golang:
	@which go > /dev/null 2>&1 || (echo "You\'re missing go executable"; @exit 1)

.PHONY: check-golangci-lint
check-golangci-lint:
	@which golangci-lint > /dev/null 2>&1 || (echo "You're missing golangci-lint executable"; exit 1)

.PHONY: install-linter
install-linter:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: fetch-lint-config
fetch-lint-config:
ifneq (,$(wildcard ${LINT_CONFIG_DIRECTORY}))
	@echo Directory is not empty. Cleaning up
	@cd ${LINT_CONFIG_DIRECTORY}
	@git pull
else
	@git clone ${LINT_CONFIG_REMOTE} ${LINT_CONFIG_DIRECTORY} > /dev/null 2>&1
endif

.PHONY: run-lint
run-lint:
	@golangci-lint run -c "${LINT_CONFIG_FILE}" ./...

.PHONY: lint
lint: check-golangci-lint run-lint

.PHONY: generate
generate:
	@go generate ./...

.PHONY: test
test:
	@go test ./...
