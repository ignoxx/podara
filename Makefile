HOST:=ec2-18-157-73-46.eu-central-1.compute.amazonaws.com
deploy: build-linux
	@scp -i "~/.ssh/podara.pem" ./bin/main-linux-arm ubuntu@$(HOST):/home/ubuntu
	@scp -R -i "~/.ssh/podara.pem" ./templates/* ubuntu@$(HOST):/home/ubuntu/templates

run:
	@env JWT_SECRET_KEY=test123 go run main.go

build-linux:
	@env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o bin/main-linux-arm main.go

build:
	@env GOARCH=arm64 GOOS=darwin go build -ldflags="-s -w" -o bin/main-darwin-arm main.go
