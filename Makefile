generate:
	go generate ./validators

migrate-up:
	go run cmd/migrate/main.go -up

migrate-down:
	go run cmd/migrate/main.go -down

build:
	go build -o ./tmp/main .

run:
	go run main.go

develop:
	air
