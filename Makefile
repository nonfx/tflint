default: build

prepare:
	go run ./plugin/stub-generator

test: prepare
	go test -timeout 5m $$(go list ./... | grep -v tflint/terraform | grep -v test-fixtures | grep -v stub-generator | grep -v integrationtest)

build:
	mkdir -p dist
	go build -v -o dist/tflint

install:
	go install

e2e: prepare install
	go test -timeout 5m ./integrationtest/inspection ./integrationtest/langserver ./integrationtest/bundled ./integrationtest/init

lint:
	go run golang.org/x/lint/golint --set_exit_status $$(go list ./... | grep -v tflint/terraform)

clean:
	rm -rf dist/

generate:
	go generate ./...

.PHONY: prepare test build install e2e lint clean generate
