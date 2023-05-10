# for fish script

build:
	env GOOS=darwin GOARCH=amd64 go build -o ./bin/mon_mac_amd64 ./cmd/simpleMonitor/main.go
	env GOOS=darwin GOARCH=arm64 go build -o ./bin/mon_mac_arm64 ./cmd/simpleMonitor/main.go
	env GOOS=linux GOARCH=amd64 go build -o ./bin/mon_linux_amd64 ./cmd/simpleMonitor/main.go