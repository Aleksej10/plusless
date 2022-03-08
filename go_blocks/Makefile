go_blocks: go_blocks.go core.go
	go build $^

clean:
	rm -f go_blocks

install: go_blocks
	mkdir -p /usr/local/bin
	cp -f go_blocks /usr/local/bin
	chmod 755 /usr/local/bin/go_blocks

uninstall:
	rm -f /usr/local/bin/go_blocks

.PHONY: clean install uninstall
