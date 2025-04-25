build:
	go build --tags "fts5" --ldflags='-w -s'.
run:
	go run --tags "fts5" .
clean:
	rm db/health_sys his
