#  Server Familiar — the checks that must pass before a phase is called done.

GO ?= go
CGO := CGO_ENABLED=0

.PHONY: check fmt vet build test agent lab demo clean

# The phase gate: run this after every task group.
check: fmt vet build test

fmt:
	@unformatted=$$(gofmt -l . 2>/dev/null); \
	if [ -n "$$unformatted" ]; then \
		echo "gofmt needed:"; echo "$$unformatted"; exit 1; \
	fi

vet:
	$(GO) vet ./...

build:
	$(CGO) $(GO) build ./...

test:
	$(GO) test ./...

# Static binaries: copy one file to a host and run it.
agent:
	$(CGO) $(GO) build -trimpath -ldflags "-s -w" -o bin/agent ./cmd/agent

lab:
	$(CGO) $(GO) build -trimpath -ldflags "-s -w" -o bin/lab ./cmd/lab

# The simulated six-familiar fleet: no servers, no accounts, no setup.
demo:
	LAB_DEMO=1 $(GO) run ./cmd/lab

clean:
	rm -rf bin dist
