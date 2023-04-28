BINARY_NAME=ops

hello:
	echo "Hello"

build:
	go build -o bin/$(BINARY_NAME) main.go

run:
	go run main.go

compile:
	echo "Compiling for every OS and Platform"
	GOOS=linux GOARCH=arm go build -o bin/$(BINARY_NAME)-linux-arm main.go
	GOOS=linux GOARCH=arm64 go build -o bin/$(BINARY_NAME)-linux-arm64 main.go
	GOOS=freebsd GOARCH=386 go build -o bin/$(BINARY_NAME)-freebsd-386 main.go

all: hello build