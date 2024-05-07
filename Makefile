.PHONY=build

build-producer:
	@go build -o bin/producer producer/main.go

run-producer: build-producer
	@./bin/producer

build-consumer:
	@go build -o bin/consumer consumer/main.go

run-consumer: build-consumer
	@./bin/consumer

test:
	@go test -v -cover ./test/...