build:
	@go build -C cmd -o ../bin/userSegment.exe

run: build
	@./bin/userSegment

test:
	@go test -v ./...