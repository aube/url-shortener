.PHONY: run
run:
	go run ./cmd/shortener/ -d=postgres://pguser:pgpass@localhost:5430/shortenerdb?sslmode=disable -c=configs -t=192.168.1.0/24

.PHONY: runfile
runfile:
	go run ./cmd/shortener/ -f=./_storage/urls_list.txt -c=configs -t=192.168.1.0/24

.PHONY: runmem
runmem:
	go run ./cmd/shortener/ -c=configs -t=192.168.1.0/24

.PHONY: build
build:
	go build -ldflags "                                     \
	-X main.buildVersion=v1.0.1                             \
	-X 'main.buildTime=$$(date +'%Y/%m/%d %H:%M:%S')'       \
	-X 'main.buildCommit=$$(git rev-parse --short HEAD)'    \
	"  -o ./cmd/shortener/shortener ./cmd/shortener/main.go \

.PHONY: buildlint
buildlint:
	go build -o ./cmd/staticlint/staticlint ./cmd/staticlint/main.go

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
	mockgen -source=./internal/app/store/interfaces.go -destination=./mocks/mocks.go

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

.PHONY: staticcheck
staticcheck:
	staticcheck ./...

.PHONY: protoc
protoc:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative internal/api/grpc/proto/urlshortener.proto

.DEFAULT_GOAL := run
