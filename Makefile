NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m


PKG_SRC := github.com/ah8ad3/gateway

.PHONY: all

all: clean deps test build


deps:
	@echo "$(OK_COLOR)==> Installing dependencies$(NO_COLOR)"
	@go get -u github.com/golang/dep/cmd/dep
	@go get -u golang.org/x/lint/golint
	@go get -u github.com/DATA-DOG/godog/cmd/godog
	@dep ensure -v -vendor-only


build:
	@echo "$(OK_COLOR)==> Building... $(NO_COLOR)"
	@/bin/sh -c "BUILD_DEFAULT=$(BUILD_DEFAULT) PKG_SRC=$(PKG_SRC) VERSION=$(VERSION) ./build/build.sh"


test: lint format vet
	@echo "$(OK_COLOR)==> Running tests$(NO_COLOR)"
	@go test -v -cover ./...


test-integration: lint format vet
	@echo "$(OK_COLOR)==> Running tests$(NO_COLOR)"
	@go test -v -cover -tags=integration ./...


format:
	@echo "$(OK_COLOR)==> checking code format with 'gofmt' tool$(NO_COLOR)"
	@gofmt -l -s cmd pkg | grep ".*\.go"; if [ "$$?" = "0" ]; then exit 1; fi


vet:
	@echo "$(OK_COLOR)==> checking code correctness with 'go vet' tool$(NO_COLOR)"
	@go vet ./...


lint:
	@echo "$(OK_COLOR)==> checking code style with 'golint' tool$(NO_COLOR)"
	@go list ./... | xargs -n 1 golint -set_exit_status


clean:
	@echo "$(OK_COLOR)==> Cleaning project$(NO_COLOR)"
	@go clean
	@rm -rf bin $GOPATH/bin

number:
	@echo "$(OK_COLOR)==> checking code numbers$(NO_COLOR)"
	@find . -name "*.go" | xargs wc -l

cov:
	@echo "$(OK_COLOR)==> Getting test coverage $(NO_COLOR)"
	export TEST=1; \
	go test ./... -coverprofile=coverage.out -cover; \
	go tool cover -html=coverage.out;
