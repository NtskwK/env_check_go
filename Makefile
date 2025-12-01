build:
	@echo building dnsbenchmark...
	go mod tidy
	GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -tags release -o release/env_checkGo-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -tags release -o release/env_checkGo-darwin-arm64 .
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -tags release -o release/env_checkGo-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -tags release -o release/env_checkGo-linux-arm64 .
	GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -tags release -o release/env_checkGo-windows-amd64.exe .
	GOOS=windows GOARCH=arm64 go build -ldflags "-s -w" -tags release -o release/env_checkGo-windows-arm64.exe .
