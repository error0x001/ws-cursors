lint_install:
	if [ ! -f $(GOPATH)/bin/golangci-lint ] ; then \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.43.0; \
    fi

lint: lint_install
	golangci-lint --version && golangci-lint run --config=.golangci.yml

gobadge_install:
	go get github.com/AlexBeauchemin/gobadge

update_coverage: gobadge_install
	go test -v ./... -covermode=count -coverprofile=coverage.out
	go tool cover -func=coverage.out -o=coverage.out
	gobadge -filename=coverage.out
