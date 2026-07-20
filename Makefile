.PHONY: build run test tidy

build:
	go build -ldflags "-X github.com/meclaw/meclaw/cmd.Version=dev" -o bin/meclaw .

run: build
	./bin/meclaw chat -c examples/config.example.json

test:
	go test ./...

tidy:
	go mod tidy
