cov:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

build:
	go build ./cmd/jbd/.

run:
	go run ./cmd/jbd/.

test:
	go test -v ./...

lint:
	golangci-lint run

fmt:
	go fmt ./...

migrations-up:
	migrate -path=./db/migrations -database=$(JUNKBOY_DB_DSN) up

migrations-down:
	migrate -path=./db/migrations -database=$(JUNKBOY_DB_DSN) down