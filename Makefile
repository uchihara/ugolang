run:
		go run main.go 'func main() int { 1+2*3; }'
test:
	go test ugolang/*.go
