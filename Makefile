code = $(CODE)
ifeq ($(code),)
	code = func main() int { 1+2*3; }
endif

run:
		go run main.go "$(code)"
test:
	go test ugolang/*.go
