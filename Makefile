.PHONY: tinypdf docker run

tinypdf:
	go build -o tinypdf -trimpath -ldflags="-s -w" main.go

docker:
	docker build -t bhupeshimself/tinypdf:latest .

run:
	docker run --rm -v $(shell pwd):/app bhupeshimself/tinypdf:latest