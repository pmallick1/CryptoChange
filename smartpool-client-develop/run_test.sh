build/env.sh go test -ldflags -s `go list ./... | grep -v vendor | grep -v experiment | grep -v ethash`
