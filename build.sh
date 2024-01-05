#!/usr/bin/env bash

# macos arm64
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o executable/darwin_arm64_inscribe main.go

sleep 3

# macos amd64
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o executable/darwin_amd64_inscribe main.go

sleep 3

# 交叉编译windows
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o executable/windows_amd64_inscribe.exe main.go

sleep 3

# 交叉编译linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o executable/linux_amd64_inscribe main.go