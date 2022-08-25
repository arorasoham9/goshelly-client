BINARY_NAME=goshelly

build:
	rm -rf bin/
	mkdir bin/
	GOOS=darwin go build -o bin/${BINARY_NAME}-darwin main.go
	GOOS=linux go build -o bin/${BINARY_NAME}-linux main.go

run:
	./bin/${BINARY_NAME}-linux

build_and_run: 
	make build 
	make run

clean:
	go clean
	rm bin/${BINARY_NAME}-darwin
	rm bin/${BINARY_NAME}-linux