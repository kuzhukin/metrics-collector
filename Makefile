ALL_TARGETS = agent server

define build
	go build -v ./cmd/$(1)
endef

all: build

.PHONY: build
build: $(patsubst %, build-%, $(ALL_TARGETS))


.PHONY: build-%
build-%:
	@echo === Building $*
	$(call build,$*)

.PHONY: test
test:
	@echo === Tests
	go test -count 1 -v ./...

clean:
	rm server && \
	rm agent

# Linter constants
LINTER := golangci-lint 

.PHONY: lint
lint:
	@echo === Lint
	$(LINTER) --version
	$(LINTER) cache clean && $(LINTER) run

genmock:
	go generate ./...