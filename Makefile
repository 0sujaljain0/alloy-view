build:
	mkdir -p bin/ && rm -rf bin/* &&  go build -o bin/main cmd/main.go 

clean_logs:
	echo "" > ./logs.log

format:
	gofmt -w .

vet:
	go vet ./...

run: format vet build clean_logs
	./bin/main

run_full: run
	go tool pprof -png cpu.prof > cpu.png && open cpu.png
