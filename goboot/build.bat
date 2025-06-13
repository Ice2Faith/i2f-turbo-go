@echo off

echo clean mod ...
go mod tidy

echo install gox ...
go install github.com/mitchellh/gox@latest

echo clean output directory ...
del /f /s /q output
rd /s /q output

echo build ...
mkdir output
:: gox -os "windows linux darwin" -arch "386 amd64 arm64" -output "output/{{.OS}}-{{.Arch}}/goboot"

gox -osarch "windows/386 windows/amd64 linux/386 linux/amd64 linux/arm64 darwin/amd64 darwin/arm64" -output "output/{{.OS}}-{{.Arch}}/goboot"
