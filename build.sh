#!/bin/bash

set -e

go build cmd/main.go
go build -o root/init cmd/init/init.go
go build -o root/init2 cmd/init/init2.go

