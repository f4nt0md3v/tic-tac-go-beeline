-include .env

build:
	@echo "  >  Building package..."
	go build -o cmd/${BIN_FILENAME} ${GO_PACKAGE_NAME}

run:
	@echo "  >  Running package..."
	go run ${GO_PACKAGE_NAME}

detect-race:
	@echo "  >  Running package in race detection mode..."
	go run -race ${GO_PACKAGE_NAME}

test:
	@echo "  >  Testing package..."
	go test ${GO_PACKAGE_NAME}

fmt:
	@echo "  >  Formating package..."
	go fmt ${GO_PACKAGE_NAME}

clean:
	@echo "  >  Cleaning up build artifacts..."
	go clean

lint:
	golangci-lint run ./...

visualize:
	go-callvis ${GO_PACKAGE_NAME}

docker-run-dev:
	docker run -it -p 9090:9090 ${DOCKER_IMAGE_NAME}

compose-up:
	docker-compose up -d

compose-down:
	docker-compose down
