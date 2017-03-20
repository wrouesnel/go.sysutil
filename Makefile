
GO_SRC := $(shell find -type f -name "*.go")
PACKAGES = $(shell go list ./... | grep -v /vendor/)

COVERAGE_DIR=.coverage
TOOLS_DIR=.tools

all: style vet test

vet:
	go vet $(PACKAGES)

# Check code conforms to go fmt
style:
	! gofmt -s -l $(GO_SRC) 2>&1 | read 2>/dev/null

# Format the code
fmt:
	gofmt -s -w $(GO_SRC)

lint: tools
	$(TOOLS_DIR)/golint $(PACKAGES)

test: tools
	@mkdir -p $(COVERAGE_DIR)
	for pkg in $(PACKAGES) ; do \
		go test -coverprofile=$(COVERAGE_DIR)/$$(echo $$pkg | tr '/' '-').test.out -covermode=count $$pkg ; \
	done
	$(TOOLS_DIR)/gocovmerge $(COVERAGE_DIR)/* > cover.out

tools: $(TOOLS_DIR)/gocovmerge $(TOOLS_DIR)/golint

$(TOOLS_DIR)/golint: 
	@mkdir -p $(TOOLS_DIR)
	go build -o $(TOOLS_DIR)/golint ./vendor/github.com/golang/lint/golint/.
	
$(TOOLS_DIR)/gocovmerge:
	@mkdir -p $(TOOLS_DIR)
	go build -o $(TOOLS_DIR)/gocovmerge ./vendor/github.com/wadey/gocovmerge/.

.PHONY: fmt style vet test tools
