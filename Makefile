PREFIX?=$(shell pwd)
EXEC_NAME=isocli

.PHONY: clean build

build: *.go
	@go build -o $(EXEC_NAME) .

clean:
	@rm $(EXEC_NAME)
