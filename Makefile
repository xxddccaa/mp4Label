.PHONY: all clean install run windows linux darwin prepare-web clean-web all-platforms

VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "v0.2.5")
PKGBASE := github.com/xd/mp4label/cmd

# Default target: build all platforms
all: all-platforms

# 准备 web 目录（用于 embed）
prepare-web:
	@echo "Preparing web directory for embedding..."
	@rm -rf cmd/mp4label/web
	@cp -r web cmd/mp4label/

# 清理 web 目录
clean-web:
	@rm -rf cmd/mp4label/web

cmd-mp4label: prepare-web
	@echo "Building mp4label version: ${VERSION}"
	CGO_ENABLED=0 go build -ldflags "-w -s -X main.version=${VERSION}" -o bin/mp4label ${PKGBASE}/mp4label

# Windows 64位版本
windows: prepare-web
	@echo "Building Windows 64-bit version: ${VERSION}"
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-w -s -X main.version=${VERSION}" -o bin/mp4label.exe ${PKGBASE}/mp4label

# Linux 版本
linux: prepare-web
	@echo "Building Linux version: ${VERSION}"
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w -s -X main.version=${VERSION}" -o bin/mp4label-linux ${PKGBASE}/mp4label

# macOS 版本
darwin: prepare-web
	@echo "Building macOS version: ${VERSION}"
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-w -s -X main.version=${VERSION}" -o bin/mp4label-darwin ${PKGBASE}/mp4label

# Build all platforms
all-platforms: windows linux darwin
	@echo "========================================="
	@echo "All platforms built successfully!"
	@echo "========================================="
	@echo "Generated binaries:"
	@ls -lh bin/
	@echo "========================================="

clean: clean-web
	rm -rf bin/

install: cmd-mp4label
	cp bin/mp4label /usr/local/bin/mp4label

run: cmd-mp4label
	./bin/mp4label web
