.PHONY: test clean

build: build-docs
	@bash -c ./build.sh

build-docs:
	@echo "Building docs..."
	@docker run --rm -v $(shell pwd):/spec redocly/cli build-docs api_v1.yml --output public/index.html --theme.openapi.disableSearch

test:
	@echo "Running tests..."
	@go test -coverprofile=coverage.out -v ./...

coverage: test
	@echo "Generating coverage report..."
	@go tool cover -html=coverage.out

clean:
	@rm -rf functions public