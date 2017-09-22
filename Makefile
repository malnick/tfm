.PHONY: all build test clean
THIS_FILE := $(lastword $(MAKEFILE_LIST))

os?="darwin"

all: clean test build-all

build-all:
	@$(MAKE) -f $(THIS_FILE) build os=linux
	@$(MAKE) -f $(THIS_FILE) build os=darwin
	@$(MAKE) -f $(THIS_FILE) build os=freebsd
	@$(MAKE) -f $(THIS_FILE) build os=windows

build:
	$(call i,building tfm for $(os))
	@bash -c "scripts/build.sh $(os)"

test:
	bash -c "./scripts/test.sh libraries unit"

docker-test:
	bash -c "./scripts/docker-test.sh"

clean:
	rm -rf ./build

# Helper Functions
define i 
	@tput setaf 6 && echo "[INFO] ==> $(1)"
	@tput sgr0
endef

define d 
	@tput setaf 7 && echo "[DEBUG] ==> $1"
	@tput sgr0
endef

define w
	@tput setaf 3 && echo "[WARN] ==> $1"
	@tput sgr0
endef

define e 
	@tput setaf 1 
	@echo "[ERROR] ==> $1"
	@tput sgr0
endef

