#!/bin/bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")" && pwd)"
APP_NAME="Stephen Bird Site.app"
APP_DIR="$ROOT/dist/$APP_NAME"
BIN_NAME="Stephen Bird Site"
SITE_DIR="$APP_DIR/Contents/Resources/site"
MACOS_DIR="$APP_DIR/Contents/MacOS"
PLIST="$APP_DIR/Contents/Info.plist"

rm -rf "$APP_DIR"
mkdir -p "$MACOS_DIR" "$SITE_DIR"

cat > "$PLIST" <<'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleExecutable</key>
	<string>Stephen Bird Site</string>
	<key>CFBundleIdentifier</key>
	<string>com.stephenmbird.site</string>
	<key>CFBundleName</key>
	<string>Stephen Bird Site</string>
	<key>CFBundlePackageType</key>
	<string>APPL</string>
	<key>CFBundleShortVersionString</key>
	<string>1.0</string>
	<key>CFBundleVersion</key>
	<string>1</string>
	<key>LSMinimumSystemVersion</key>
	<string>11.0</string>
</dict>
</plist>
EOF

rsync -a --delete \
	--exclude 'dist/' \
	--exclude 'tools/' \
	--exclude 'build-macos-app.sh' \
	--exclude '.DS_Store' \
	"$ROOT"/ "$SITE_DIR"/

GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o "$MACOS_DIR/$BIN_NAME" ./tools/site-launcher
chmod +x "$MACOS_DIR/$BIN_NAME"

echo "Built: $APP_DIR"
