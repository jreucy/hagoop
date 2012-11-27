all: go/src/mrlib/mrlib.go worker request server main

main: go/src/client/client-impl.go go/src/main/main.go
	go build go/src/main/main.go

worker: go/src/worker/worker.go
	go build go/src/worker/worker.go

request: go/src/request/request.go
	go build go/src/request/request.go

server: go/src/server/server.go
	go build go/src/server/server.go

clean:
	rm main worker request server
