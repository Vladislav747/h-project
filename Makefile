build:
	@go build -o cmd

run:
	@go run ./cmd/h-project

seed:
	@go run scripts/seed.go

test:
	@go test -v ./...