
run:
	go run ./cmd/jbd/.

test:
	go test -v ./...

lint:
	golangci-lint run --enable-all

migrations-up:
	migrate -path=./db/migrations -database=$(JUNKBOY_DB_DSN) up

migrations-down:
	migrate -path=./db/migrations -database=$(JUNKBOY_DB_DSN) down