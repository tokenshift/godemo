#!/bin/sh

APP_NAME='godemo'
APP_REPO="github.com/tokenshift/godemo"

IFS='/'

set -x

go tool dist list | while read os arch; do
	env GOOS=$os GOARCH=$arch go build -o "target/${APP_NAME}.${os}_${arch}" "$APP_REPO"
done
