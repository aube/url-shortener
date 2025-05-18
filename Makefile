.PHONY: run
run:
	go run ./cmd/shortener/ -d=postgres://videouser:videopass@localhost:5432/videodb?sslmode=disable

.PHONY: build
build:
	go build -o ./cmd/shortener *.go

.PHONY: profbase
profbase:
	curl -sK -v http://localhost:8080/debug/pprof/heap?seconds=10 > ./profiles/base.pprof

.PHONY: profres
profres:
	curl -sK -v http://localhost:8080/debug/pprof/heap?seconds=10 > ./profiles/result.pprof

.PHONY: profdiff
profdiff:
	go tool pprof -top -diff_base=profiles/base.pprof profiles/result.pprof

.PHONY: mocks
mocks:
	mockgen -source=./internal/api/interfaces.go -destination=./mocks/api/mocks.go

.PHONY: test
test:
	go test -v -timeout 30s ./...

.PHONY: race
race:
	go test -v -race -timeout 30s ./...

.PHONY: cover
cover:
	go test -v -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o test-coverage.html
	rm coverage.out

.DEFAULT_GOAL := run
