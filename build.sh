#!/bin/sh

GOOS=darwin  GOARCH=amd64 go build -o ./hostsmerge_darwin_amd64.exe  github.com/apaxa-io/hostsmerge
GOOS=linux   GOARCH=386   go build -o ./hostsmerge_linux_386.exe     github.com/apaxa-io/hostsmerge
GOOS=linux   GOARCH=amd64 go build -o ./hostsmerge_linux_amd64.exe   github.com/apaxa-io/hostsmerge
GOOS=windows GOARCH=386   go build -o ./hostsmerge_windows_386.exe   github.com/apaxa-io/hostsmerge
GOOS=windows GOARCH=amd64 go build -o ./hostsmerge_windows_amd64.exe github.com/apaxa-io/hostsmerge
