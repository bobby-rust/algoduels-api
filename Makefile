build:
	@go build -o bin/test-rest-api

run: build	
	@./bin/test-rest-api

test:
	@go test -v ./...