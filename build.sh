#!/bin/bash

BINARY=trm

cd $(dirname "$0")
go build -o $BINARY main.go 
if [ -d "bin" ]; then
    mv $BINARY bin/
else
    mkdir bin
    mv $BINARY bin/
fi