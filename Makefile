BUILD   := target/
OBJECT  := blap
TARGET  := $(BUILD)$(OBJECT)
DESTDIR := /usr/local/bin
GOFLAGS := -trimpath -buildmode=pie -mod=readonly -modcacherw -buildvcs=false
VERSION ?= $(shell git log -n 1 --format=%h)

all: $(TARGET)

build: $(TARGET)

$(TARGET): cmd/main.go $(shell find internal/ -type f | grep -v testdata)  go.*
	go build $(GOFLAGS) -ldflags "$(LDFLAGS) -X main.version=$(VERSION)" -o $@ cmd/main.go

install:
	install -m755 $(TARGET) $(DESTDIR)/$(OBJECT)

clean:
	@rm -rf $(BUILD)

unittest:
	go test ./...

check: $(TARGET) unittest
