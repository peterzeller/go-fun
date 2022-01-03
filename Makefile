

build:
	go1.18beta1 build ./...

test:
	go1.18beta1 test -v -coverprofile=profile.cov ./...
	go tool cover -html=profile.cov -o cover.html

clean:
	find . -type f -name '*.fail' -delete
	rm -f cover.html
	rm -f profile.cov