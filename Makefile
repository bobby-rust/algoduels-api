
build:
	@cd src && go build -o ../bin/main

run: build
	@./bin/main

test:
	@cd src && go test -v ./...
