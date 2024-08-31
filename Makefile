build:
	go build -race -o ./bin/main ./cmd/main/

run:
	go run -race ./cmd/main/

display:
	go run ./cmd/display/
