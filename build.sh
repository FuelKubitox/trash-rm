#!/bin/bash

cd $(dirname "$0")
go build main.go
if [ -d "bin" ]; then
    mv main bin/
else
    mkdir bin
    mv main bin/
fi