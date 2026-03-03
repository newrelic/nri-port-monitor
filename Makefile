MAIN=src/port-monitor.go
BIN=bin/
OUT=nri-port-monitor

all: macos-arm macos-intel linux-arm linux-intel windows

macos-arm:
	GOOS=darwin GOARCH=arm64 go build -o $(BIN)/macos-arm64/$(OUT) $(MAIN)

macos-intel:
	GOOS=darwin GOARCH=amd64 go build -o $(BIN)/macos-amd64/$(OUT) $(MAIN)

linux-arm:
	GOOS=linux GOARCH=arm64 go build -o $(BIN)/linux-arm64/$(OUT) $(MAIN)

linux-intel:
	GOOS=linux GOARCH=amd64 go build -o $(BIN)/linux-amd64/$(OUT) $(MAIN)

windows:
	GOOS=windows GOARCH=amd64 go build -o $(BIN)/windows-amd64/$(OUT).exe $(MAIN)

test:
	go test -v ./src/

clean:
	rm -rf $(BIN)
