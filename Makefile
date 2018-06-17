commit = $(shell git log --pretty=format:'%h' -n 1)
now = $(shell date "+%Y-%m-%d %T UTC%z")
compiler = $(shell go version)
branch = $(shell git rev-parse --abbrev-ref HEAD)

IMAGE_NAME := earvin/rider:${commit}

all: test build image

test:
	@echo "Running tests"
	@docker-compose -f docker-compose.test.yml up 	\
	 		--build 								\
			--abort-on-container-exit				\
			--force-recreate 						\
			--quiet-pull							\
			--no-color								\
			--remove-orphans 						\
			--timeout 20
	@docker-compose rm -f

build:
	@echo "Compiling the binary"
	@CGO_ENABLED=0  GOBIN=$(PWD)/bin go install  -v	\
	    -ldflags                              		\
	       "-X 'main.branch=$(branch)'        		\
	       	-X 'main.sha=$(commit)'           		\
	       	-X 'main.compiledAt=$(now)'       		\
	       	-X 'main.compiler=$(compiler)'			\
			-s -w"   								\
	    -a -installsuffix cgo ./...


image:
	@echo "Building Docker Image"
	@(docker build -t $(IMAGE_NAME) .)
