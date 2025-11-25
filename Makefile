build:
	mkdir -p bin/ && rm -rf bin/* &&  go build -o bin/main cmd/main.go 


run: build
	./bin/main
