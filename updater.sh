#!/bin/bash
echo "$1"
git tag "$1" && git push --tags
go-bindata -pkg main -o bindata.go -prefix "data/" data/GeoLite2-City.mmdb
gox -output="release/{{.Dir}}_{{.OS}}_{{.Arch}}" -osarch="darwin/386 darwin/amd64 linux/386 linux/amd64 linux/arm windows/386 windows/amd64"
rm bindata.go
github-release release --user ad --repo gogeo --tag "$1" --name "$1"
github-release upload --user ad --repo gogeo --tag "$1" --name "gogeo_darwin_386" --file release/gogeo_darwin_386
github-release upload --user ad --repo gogeo --tag "$1" --name "gogeo_darwin_amd64" --file release/gogeo_darwin_amd64
github-release upload --user ad --repo gogeo --tag "$1" --name "gogeo_linux_386" --file release/gogeo_linux_386
github-release upload --user ad --repo gogeo --tag "$1" --name "gogeo_linux_amd64" --file release/gogeo_linux_amd64
github-release upload --user ad --repo gogeo --tag "$1" --name "gogeo_linux_arm" --file release/gogeo_linux_arm
github-release upload --user ad --repo gogeo --tag "$1" --name "gogeo_windows_386.exe" --file release/gogeo_windows_386.exe
github-release upload --user ad --repo gogeo --tag "$1" --name "gogeo_windows_amd64.exe" --file release/gogeo_windows_amd64.exe
