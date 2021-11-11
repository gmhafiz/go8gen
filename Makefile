build:
	go build -ldflags="-w -s" -o bin/go8 cmd/go8/main.go
	mv bin/go8 ~/go/bin/go8
