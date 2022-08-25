BINARY_NAME=goshelly

build:
	rm -rf bin/
	mkdir bin
	GOOS=darwin go build -o ${BINARY_NAME}-darwin main.go
	GOOS=linux go build -o /${BINARY_NAME}-linux main.go

run:
	./${BINARY_NAME}-linux

build_and_run: 
	build run

clean:
	go clean
	rm ${BINARY_NAME}-darwin
	rm ${BINARY_NAME}-linux