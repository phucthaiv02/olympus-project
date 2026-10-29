PROBLEM_ID := kh7d8qtdcfr0k5mt1rvftcvebn8gb2va
VERSION := 0
REPO_URL := https://github.com/davidhalter/jedi
COMMIT := 2af1a0f3943f801f5de755f8ffd251f0385a3607
IMAGE_TAG := olympus-$(PROBLEM_ID)-v$(VERSION)
REPO_DIR := jedi

.PHONY: help clone checkout apply-test-patch apply-solution-patch build-image

help:
	@echo "clone                Clone the source repo recursively and checkout the challenge commit"
	@echo "checkout             Re-checkout the challenge commit inside ./repo"
	@echo "apply-test-patch     Apply test.patch onto ./repo if present"
	@echo "apply-solution-patch Apply solution.patch onto ./repo if present"
	@echo "build-image          Build from the checked-out repo context on linux/amd64"

clone:
	@if [ -z "$(REPO_URL)" ]; then echo "REPO_URL is empty"; exit 1; fi
	@if [ -z "$(COMMIT)" ]; then echo "COMMIT is empty"; exit 1; fi
	@if [ ! -d "$(REPO_DIR)/.git" ]; then git clone --recursive "$(REPO_URL)" "$(REPO_DIR)"; else echo "repo already exists at $(REPO_DIR)"; fi
	@cd "$(REPO_DIR)" && git fetch --all --tags && git checkout "$(COMMIT)" && git submodule update --init --recursive

checkout:
	@if [ -z "$(COMMIT)" ]; then echo "COMMIT is empty"; exit 1; fi
	@if [ ! -d "$(REPO_DIR)/.git" ]; then echo "repo not present at $(REPO_DIR); run make clone"; exit 1; fi
	@cd "$(REPO_DIR)" && git fetch --all --tags && git checkout "$(COMMIT)" && git submodule update --init --recursive

apply-test-patch: checkout
	@if [ -f test.patch ]; then cd "$(REPO_DIR)" && git apply ../test.patch; else echo "test.patch not present"; fi

apply-solution-patch: checkout
	@if [ -f solution.patch ]; then cd "$(REPO_DIR)" && git apply ../solution.patch; else echo "solution.patch not present"; fi

build-image: checkout
	@if [ ! -f Dockerfile ]; then echo "Dockerfile not present"; exit 1; fi
	@cd "$(REPO_DIR)" && rm -f .dockerignore && docker build --platform linux/amd64 -f ../Dockerfile -t "$(IMAGE_TAG)" .
