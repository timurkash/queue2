
build:
	mkdir -p bin && \
	GO111MODULE=on CGO_ENABLE=0 GOOS=linux GOARCH=amd64 \
	go build -o ./bin/ ./...

test:
	go test ./internal/test/...
