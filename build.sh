#! /bin/sh
go build -ldflags "-X 'github.com/clcert/osr/cmd.Build=$(date -u '+%Y-%m-%d %H:%M:%S')'"
