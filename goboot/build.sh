#!bin/bash

echo clean mod ...
go mod tidy

echo install gox ...
go install github.com/mitchellh/gox@latest

echo clean output directory ...
rm -rf output

echo build ...
mkdir output
# gox -parallel=-1 -os "windows linux darwin" -arch "386 amd64 arm64" -output "output/{{.OS}}-{{.Arch}}/goboot"
# -parallel=-1        Amount of parallelism, defaults to number of CPUs

gox -parallel=2 -osarch "windows/386 windows/amd64 linux/386 linux/amd64 linux/arm64 darwin/amd64 darwin/arm64" -output "output/{{.OS}}-{{.Arch}}/goboot"
