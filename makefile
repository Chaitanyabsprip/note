.DEFAULT_GOAL:=./bin/note
INSTALL_PATH=/usr/local/bin/note

./bin/note: main.go
	@go build -o ./bin/note

clean:
	@rm -rd ./bin

install: ./bin/note
	@install ./bin/pomo ${INSTALL_PATH}

uninstall: clean
	@rm ${INSTALL_PATH}
