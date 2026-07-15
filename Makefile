APP := originalrush
DIST_DIR := dist
PACKAGE_NAME := DiamondRush-macos-arm64
PACKAGE_DIR := $(DIST_DIR)/$(PACKAGE_NAME)
PACKAGE_ZIP := $(DIST_DIR)/$(PACKAGE_NAME).zip

.PHONY: package package-macos-arm64 clean-package

package: package-macos-arm64

package-macos-arm64:
	rm -rf "$(PACKAGE_DIR)" "$(PACKAGE_ZIP)"
	mkdir -p "$(PACKAGE_DIR)"
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -trimpath -o "$(PACKAGE_DIR)/$(APP)" ./cmd/originalrush
	cp -R decoded "$(PACKAGE_DIR)/decoded"
	ditto -c -k --keepParent "$(PACKAGE_DIR)" "$(PACKAGE_ZIP)"
	@echo "Package created: $(PACKAGE_ZIP)"

clean-package:
	rm -rf "$(DIST_DIR)"
