VERSION = 0.1.0
PACKAGE = github.com/KyleBanks/readme
RELEASE_PLATFORMS = darwin/386 darwin/amd64 linux/386 linux/amd64 linux/arm windows/386 windows/amd64

example: | install
	@readme KyleBanks/depth
.PHONY: example

install:
	@go install -v $(PACKAGE)
	@echo "Installed."
.PHONY: install

release:
	@gox -osarch="$(RELEASE_PLATFORMS)" \
        -output "bin/{{.Dir}}_$(VERSION)_{{.OS}}_{{.Arch}}" $(PACKAGE)
.PHONY: release
