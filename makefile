all: go/src/mrlib/mrlib.go go/bin/worker go/bin/request go/bin/server go/bin/wc

go/bin/wc: go/src/client/client-impl.go go/src/main/main.go
	go build -o go/bin/wc go/src/main/main.go

go/bin/worker: go/src/worker/worker.go
	go build -o go/bin/worker go/src/worker/worker.go

go/bin/request: go/src/request/request.go
	go build -o go/bin/request go/src/request/request.go

go/bin/server: go/src/server/server.go
	go build -o go/bin/server go/src/server/server.go

clean:
	rm go/bin/wc go/bin/worker go/bin/request go/bin/server
