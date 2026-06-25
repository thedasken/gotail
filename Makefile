run:
	@go run .

test:
	@go test -v ./...

build:
	@go build .

clean:
	@rm -f gotail

.PHONY: run test build clean