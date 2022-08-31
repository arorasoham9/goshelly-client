BINARY_NAME=goshelly

build:
	rm -rf bin/
	mkdir bin/
	GOOS=darwin go build -o bin/${BINARY_NAME}_darwin main.go
	GOOS=linux go build -o bin/${BINARY_NAME}_linux main.go
run:
	./bin/${BINARY_NAME}_linux assess

build_and_run: 
	make build 
	make run 

clean:
	go clean
	rm bin/${BINARY_NAME}_darwin
	rm bin/${BINARY_NAME}_linux