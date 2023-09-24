file=$(date '+%Y%m%d%H%M%S')

echo "Building Linux versions..."
GOOS=linux GOARCH=amd64 go build -o build/$file/Anonymous-Bot-64 -trimpath -ldflags "-s -w -buildid=" cmd/main.go
GOOS=linux GOARCH=arm64 go build -o build/$file/Anonymous-Bot-arm64 -trimpath -ldflags "-s -w -buildid=" cmd/main.go

# 打包 Windows 版本
echo "Building Windows versions..."
GOOS=windows GOARCH=amd64 go build -o build/$file/Anonymous-Bot-windows.exe -trimpath -ldflags "-s -w -buildid=" cmd/main.go
GOOS=windows GOARCH=arm64 go build -o build/$file/Anonymous-Bot-windows-arm64.exe -trimpath -ldflags "-s -w -buildid=" cmd/main.go

# 打包 macOS 版本
echo "Building macOS versions..."
GOOS=darwin GOARCH=arm64 go build -o build/$file/Anonymous-Bot-macos-arm64 -trimpath -ldflags "-s -w -buildid=" cmd/main.go

echo "打包完成"