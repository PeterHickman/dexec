#!/bin/sh

BINARY='/usr/local/bin'

echo "Building dexec"
go build dexec.go

echo "Installing dexec to $BINARY"
install -v dexec $BINARY

echo "Removing the build"
rm dexec
