GOCMD = go
CGO_ENABLED = 0
GOGET = $(GOCMD) get -d -v
GOBUILD = CGO_ENABLED=$(CGO_ENABLED) $(GOCMD) build -ldflags="-s -w" -gcflags="-trimpath=$(GOPATH)" -asmflags="-trimpath=$(GOPATH)" -a -installsuffix cgo
GOBUILD_DOCKER = CGO_ENABLED=$(CGO_ENABLED) $(GOCMD) build -ldflags="-s -w" -a -installsuffix cgo

all: build
build:
	for f in analyzers/**/*.go; \
		do $(GOGET) "./$${f%/*}" && \
		$(GOBUILD) -o "$${f%%.*}" $$f; \
	done
docker-build:
	sudo docker run -it --rm -v "$(PWD)":/data -w /data golang:latest /bin/bash -c \
		'for f in analyzers/**/*.go; \
			do $(GOGET) "./$${f%/*}" && \
			$(GOBUILD_DOCKER) -o "$${f%%.*}" $$f; \
		done'
