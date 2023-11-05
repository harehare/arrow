main := "."
module := "github.com/harehare/fcd"
lint := "staticcheck"
sec := "gosec"
target := "./..."

setup:
  go get honnef.co/go/tools/cmd/staticcheck
  go get github.com/securego/gosec/v2/cmd/gosec
  go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest
  go install github.com/cosmtrek/air@latest

run *args:
	go run {{ main }} {{args}}

build:
	go build -o dist/arrow {{ main }}

watch:
	air -c .air.toml

test:
	go test {{ target }}

lint:
	{{ lint }} {{ module }}
	{{ sec }} {{ target }}
