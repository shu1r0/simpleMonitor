build:
	go build -o main_mac main.go
	GOOS=linux GOARCH=amd64 go build -o main_linux main.go