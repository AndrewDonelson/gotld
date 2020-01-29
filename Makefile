
phony: test bench build run

test:
	go test -v .

bench:
	go test -bench .

build:
	go build -o ./example/main.go

run:
	./example/example
