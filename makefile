-include .env

.PHONY: run test build docker-build docker-run

run:
	go run ./cmd

test:
	GOCACHE=/private/tmp/ticktick-ai-go-cache go test ./...

build:
	go build -o bin/ticktick-ai ./cmd

docker-build:
	docker build -f dockerfile -t ticktick-ai:latest .

docker-run:
	docker run --env-file .env ticktick-ai:latest
