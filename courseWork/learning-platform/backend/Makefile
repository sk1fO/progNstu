.PHONY: build-docker run build

build-docker:
	docker build -f docker/Dockerfile.cpp -t cpp-runner .

run:
	go run main.go

build:
	go build -o platform

clean:
	rm -f platform