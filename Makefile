
check:
	go vet ./...
	go install honnef.co/go/tools/cmd/staticcheck@v0.3.2
	staticcheck ./...

build:
	go build ./...

test:
	go test -v -coverprofile=profile.cov ./...
	go tool cover -html=profile.cov -o cover.html

clean:
	find . -type f -name '*.fail' -delete
	rm -f cover.html
	rm -f profile.cov