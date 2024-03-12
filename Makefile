ALL_TARGETS = agent server

define build
	go build -o ./cmd/$(1)/ -v ./cmd/$(1)/ 
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
	go test -count 1 -v -cover ./...

define clean
	rm ./cmd/$(1)/$(1)
endef

.PHONY: clean
clean: $(patsubst %, clean-%, $(ALL_TARGETS))

clean-%:
	@echo === Cleaning $*
	$(call clean,$*)

# Linter constants
LINTER := golangci-lint 

.PHONY: lint
lint:
	@echo === Lint
	$(LINTER) --version
	$(LINTER) cache clean && $(LINTER) run

generate:
	go generate ./...

.PHONY: test_with_coverage
test_with_coverage:
	go test -count=1 -coverprofile=coverage.out ./...

.PHONY: coverage_total
coverage_total: test_with_coverage
	PERCENTAGE=$$(go tool cover -func=coverage.out | grep "total:" | tr -s '\t' | cut -f3); \
	echo Total coverage: $${PERCENTAGE}

.PHONY: coverage_html
coverage_html: test_with_coverage
	go tool cover -html=coverage.out

.PHONY: server_benchmark_cpu_profile
server_benchmark_cpu_profile:
	go test ${PWD}/internal/server -bench BenchmarkJSONRouter -cpuprofile=profiles/server_cpu.pprof 
	
server_cpu_profile: server_benchmark_cpu_profile
	go tool pprof -http=:8282 server.test profiles/server_cpu.pprof

.PHONY: server_benchmark_mem_profile
server_benchmark_mem_profile:
	go test ${PWD}/internal/server -bench BenchmarkJSONRouter -memprofile=profiles/server_mem.pprof 
	
server_mem_profile: server_benchmark_mem_profile
	go tool pprof -http=:8282 server.test profiles/server_mem.pprof

run_staticlint: build-staticlint
	cmd/staticlint/staticlint ./...

gen_proto:
	protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  ./internal/proto/metric.proto
