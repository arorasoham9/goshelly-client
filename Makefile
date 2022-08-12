BINARY_NAME=goshelly

build:
	GOOS=darwin go build -o ${BINARY_NAME}-Darwin main.go
	GOOS=linux go build -o ${BINARY_NAME}-Linux main.go

run:
	./${BINARY_NAME}-${uname}

build_and_run: 
	build run

clean:
	go clean
	rm ${BINARY_NAME}-darwin
	rm ${BINARY_NAME}-linux