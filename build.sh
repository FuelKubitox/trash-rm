#!/bin/bash

go build main.go
if [ -d "bin" ]; then
    mv main bin/
else
    mkdir bin
    mv main bin/
fi