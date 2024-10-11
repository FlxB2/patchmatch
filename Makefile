TAG='patchmatch_test'
BINARY_NAME=patchmatch
INSTALL_PATH=/usr/local/bin

all: build

build:
	go build -o $(BINARY_NAME)

install: build
	install -m 755 $(BINARY_NAME) $(INSTALL_PATH)

clean:
	rm -f $(BINARY_NAME)

uninstall:
	rm -f $(INSTALL_PATH)/$(BINARY_NAME)

tests:
	docker build . --file test/Dockerfile --tag ${TAG}
	docker run -it ${TAG}
	docker image rm -f ${tag}