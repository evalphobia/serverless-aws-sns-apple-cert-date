.PHONY: init build deploy

LINT_OPT := -E gofmt \
            -E golint \
			-E gosec \
			-E misspell \
			-E whitespace \
			-E stylecheck


init:
	@type sls > /dev/null || npm install -g serverless@1.63.0
	go get -v ./...

lint:
	@type golangci-lint > /dev/null || go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint $(LINT_OPT) run ./...


build:
	GOOS=linux go build -o bin/serverless ./

deploy: build
	sls deploy -v
