#!/bin/bash

set -e

go build cmd/main.go
go build -o root/init cmd/init/init.go
gcc -o root/main hack/main.c
