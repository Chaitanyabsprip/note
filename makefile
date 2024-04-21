.DEFAULT_GOAL:=build
SOURCES := $(shell find . -type f -name '*.go')

build: cmd/note/main.go
	@mkdir bin 2> /dev/null
	@go build -o ./bin ./cmd/note

note:
	@mkdir bin 2> /dev/null
	@go build -o ./bin ./cmd/note

clean:
	@rm -rd ./bin

install:
	@go install ./cmd/note

uninstall: clean
	@rm "$(which note)" "$(which notes)"
