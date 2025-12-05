build:
	mkdir -p tmp/ && rm -rf tmp/* &&  go build -o tmp/main cmd/main.go 

clean_logs:
	echo "" > ./logs.log

format:
	gofmt -w .
	
generate_templates:
	templ generate

run: format generate_templates build clean_logs
	./tmp/main

run_full: run
	go tool pprof -png cpu.prof > cpu.png && open cpu.png

tailwatch:
	npx tailwindcss -i ./static/css/input.css -o ./static/css/tailwind.css --watch
