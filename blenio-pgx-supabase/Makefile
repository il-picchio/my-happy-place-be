# Makefile
ENV_LOAD := set -a && source .env && set +a

dev:
	$(ENV_LOAD) && air

build:
	@mkdir -p tmp
	@go build -o tmp/app .

run:
	@$(ENV_LOAD) && ./tmp/app

test:
	go test ./...

tidy:
	go mod tidy
