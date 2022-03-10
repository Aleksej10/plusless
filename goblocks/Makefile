goblocks: goblocks.go core.go
	go build $^

clean:
	rm -f goblocks

install: goblocks
	mkdir -p /usr/local/bin
	cp -f goblocks /usr/local/bin
	chmod 755 /usr/local/bin/goblocks

uninstall:
	rm -f /usr/local/bin/goblocks

.PHONY: clean install uninstall
