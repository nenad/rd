all: deps test

deps:
	@echo "Installing dependencies..."
	go get -u github.com/stretchr/testify

test:
	go test -v

