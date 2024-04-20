.DEFAULT_GOAL:=build
SOURCES := $(shell find . -type f -name '*.go')

build: cmd/note/main.go cmd/notes/main.go
	@mkdir bin 2> /dev/null
	@go build -o ./bin ./cmd/note ./cmd/notes

note:
	@mkdir bin 2> /dev/null
	@go build -o ./bin ./cmd/note

notes:
	@mkdir bin 2> /dev/null
	@go build -o ./bin ./cmd/notes

clean:
	@rm -rd ./bin

install:
	@go install ./cmd/note ./cmd/notes

uninstall: clean
	@rm "$(which note)" "$(which notes)"
