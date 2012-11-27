all: worker request server main

main:
	go build go/src/main/main.go

worker:
	go build go/src/worker/worker.go

request:
	go build go/src/request/request.go

server:
	go build go/src/server/server.go

clean:
	rm main worker request server
