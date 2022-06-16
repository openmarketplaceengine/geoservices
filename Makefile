
PACKAGES = $(shell go list -f '{{.Dir}}/...' -m)

test:
	@for package in $(PACKAGES); do \
		go test $$package -test.v || exit 1; \
	done

lint:
	golangci-lint run
