
.PHONY: build.docker
build.docker:
	docker run -it -w /tmp/go -v $(shell pwd):/tmp/go tinygo/tinygo:latest /bin/bash -c \
		'tinygo build -o ./${name}/main.go.wasm -scheduler=none -target wasi -wasm-abi=generic ./${name}/main.go'

.PHONY: test
test:
	go test -v -tags=proxytest ./${name}/...

.PHONY: run
run:
	docker-compose -f ${name}/docker-compose.yaml up

.PHONY: stop
stop:
	docker-compose -f ${name}/docker-compose.yaml down
