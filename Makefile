all:

server:
	go build -o ./bin/server ./cmd/server/*.go

client:
	go build -o ./bin/client ./cmd/client/*.go

stress_test:
	go build -o ./bin/stress_test ./tools/stress_test/*.go

run_server: server
	./bin/server

run_client: client
	./bin/client

run_stress_test: stress_test
	./bin/stress_test

pprof:
	curl localhost:8080/debug/pprof/heap --output ./bin/heap
	pprof -http=":8081" ./bin/heap

clean:
	@echo "Cleaning up..."
	rm ./bin/*
