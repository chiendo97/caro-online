all:

server:
	go build -o ./bin/server ./cmd/server/*.go

client:
	go build -o ./bin/client ./cmd/client/*.go

stress_test:
	go build -o ./bin/stress_test ./tools/stress_test/*.go

run_server:
	./bin/server

run_client:
	./bin/client

run_stress_test:
	./bin/stress_test

pprof:
	go tool pprof -http ':8081' http://localhost:8080/debug/pprof/profile

clean:
	@echo "Cleaning up..."
	rm ./bin/*
