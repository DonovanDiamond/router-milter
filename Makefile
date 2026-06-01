VERSION ?= v0.2.0
COMMIT := $(shell git rev-parse --short HEAD)

.PHONY: build
build:
	@echo "Building with version: $(VERSION) and commit: $(COMMIT)"
	@go build -ldflags="-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT)"

.PHONY: docker
docker: build
	@echo "Building Docker image with version: $(VERSION) and commit: $(COMMIT)"
	@sudo docker build --build-arg VERSION=$(VERSION) --build-arg COMMIT=$(COMMIT) \
		-t router-milter:$(VERSION) .
	@sudo docker save router-milter:$(VERSION) -o router-milter-$(VERSION).tar
	@sudo chown $(USER):$(USER) router-milter-$(VERSION).tar

