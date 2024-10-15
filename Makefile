clean:
	rm -rf ./bin
	rm coverage.out

build:
	go build -o bin/gq ./

clean.build: clean build

run:
	go run ./...

fmt:
	gofmt -w ./

make doc:
	godoc -http=:6060

test:
	go clean -testcache
	go test ./... -cover -coverprofile=coverage.out

test.v:
	go clean -testcache
	go test ./... -cover -coverprofile=coverage.out -v

test.cov:
	go tool cover -html=coverage.out

test.bench:
	go clean -testcache
	go test -bench=. -benchmem

test.bench.profile:
	go clean -testcache
	go test -bench=. -benchmem -memprofile memory.out

.PHONY: test
