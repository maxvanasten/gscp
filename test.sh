#!/usr/bin/env bash
go clean -testcache
{
	printf '# Latest Test Results\n\n```text\n'
	go test -v ./...
	printf '```\n'
} > TESTS.md

bat TESTS.md
