.PHONY=lint
lint: tidy
	@.bin/golangci-lint run --config=.golangci-lint.yml ./...

.PHONY=lint-fix
lint-fix: tidy
	@.bin/golangci-lint run --config=.golangci-lint.yml --fix ./...

.PHONY=test
test: tidy
	@go test ./...

.PHONY=install-lint
install-lint:
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b .bin v2.2.2

.PHONY=tidy
tidy:
	@go mod tidy

.PHONY=mocks
mocks: install-mockery
	@mockery --all --config .mockery.yaml

.PHONY=clean-mocks
clean-mocks:
	@find . -name '*_mock.go' -delete

.PHONY=install-mockery
install-mockery:
	@go install github.com/vektra/mockery/v2@latest

.PHONY=autogen
autogen: mocks

.PHONY=clean-autogen
clean-autogen: clean-mocks

.PHONY=up
up:
	@docker compose -f docker-compose-local.yaml up -d

.PHONY=exec
exec:
	@docker compose -f docker-compose-local.yaml exec dev bash

.PHONY=down
down:
	@docker compose -f docker-compose-local.yaml down
