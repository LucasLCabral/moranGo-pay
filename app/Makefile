.PHONY: build build-linux clean

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o bootstrap cmd/api/main.go
	chmod +x bootstrap
	zip -j lambda.zip bootstrap

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bootstrap cmd/api/main.go
	chmod +x bootstrap
	zip -j lambda.zip bootstrap

clean:
	rm -f bootstrap lambda.zip
