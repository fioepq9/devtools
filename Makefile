.PHONY: test
test:
	[[ -d tmp ]] || mkdir -p tmp
	go test -covermode=count -coverprofile=tmp/coverprofile.out ./...
	go tool cover -func=tmp/coverprofile.out

.PHONY: html
html: test
	go tool cover -html=tmp/coverprofile.out