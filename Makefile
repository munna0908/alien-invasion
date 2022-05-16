.PHONY: mod
mod:
	go mod tidy
.PHONY: lint
lint:
	golangci-lint run -E whitespace -E wsl -E wastedassign -E unconvert -E tparallel -E thelper -E stylecheck -E prealloc \
	-E predeclared -E nlreturn -E misspell -E makezero -E lll -E importas -E ifshort -E gosec -E  gofmt -E goconst \
	-E forcetypeassert -E dogsled -E dupl -E errname -E errorlint -E nolintlint

.PHONY: install
install:
	go install ./cmd/alieninvasion
.PHONY: build
build:
	go build -o alieninvasion ./cmd/alieninvasion/main.go
.PHONY: test
test:
	go test -v -count=1 ./simulation