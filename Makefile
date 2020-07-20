-include .env

build:
	@echo "  >  Building package..."
	go get
	cd web; yarn install
	cd web; yarn build
	docker-compose up -d
	go build -o ${BIN_FILENAME} ${GO_PACKAGE_NAME}

run:
	@echo "  >  Running package..."
	./${BIN_FILENAME}
