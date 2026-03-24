build:
	go build .

run:
	go run rest-api

clean:
	rm rest-api

test:
	go test ./...