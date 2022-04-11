NAME = rblocks

release:
	cargo build --release

install: release
	mkdir -p /usr/local/bin
	cp -f target/release/$(NAME) /usr/local/bin/$(NAME)
	chmod 755 /usr/local/bin/$(NAME)

uninstall:
	rm -f /usr/local/bin/$(NAME)

.PHONY: clean install uninstall release
