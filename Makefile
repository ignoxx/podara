run:
	@env JWT_SECRET_KEY=test123 go run main.go

build:
	@env GOARCH=arm64 GOOS=darwin go build -ldflags="-s -w" -o bin/main main.go
