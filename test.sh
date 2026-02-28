#!/usr/bin/env bash
go test -v ./... > TESTS.md
bat TESTS.md
