.PHONY: run build clean tidy docker-build docker-run

# Load .env if present
ifneq (,$(wildcard .env))
	include .env
	export
endif

run:
	@go run .

build:
	@echo "Building..."
	@CGO_ENABLED=1 go build -ldflags="-w -s" -o bot .
	@echo "✅ Built: ./bot"

clean:
	@rm -f bot
	@rm -rf data/
	@echo "✅ Cleaned"

tidy:
	@go mod tidy

test:
	@go test ./... -v

docker-build:
	@docker build -t whatsapp-bot .

docker-run:
	@docker run -it \
		-v $(PWD)/data:/app/data \
		-e CONNECT_METHOD=$(CONNECT_METHOD) \
		-e PHONE_NUMBER=$(PHONE_NUMBER) \
		-e BOT_PREFIX=$(BOT_PREFIX) \
		-e BOT_NAME=$(BOT_NAME) \
		-e OWNER_JID=$(OWNER_JID) \
		whatsapp-bot
