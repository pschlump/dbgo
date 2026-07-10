#!/bin/bash
#
# Integration check for ChkEnv: verifies that an environment variable set in
# the parent process is read correctly by a separate program built from ./test.
# The same behavior is also covered in-process by TestChkEnvReadsRealEnv in
# env_test.go.
#
set -e

tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT
bin="$tmpdir/chkenv"

go build -o "$bin" ./test

export YepYep=Yes
got="$("$bin")"
if [ "$got" != "true" ]; then
	echo "FAIL: YepYep=Yes expected 'true', got '$got'"
	exit 1
fi

export YepYep=No
got="$("$bin")"
if [ "$got" != "false" ]; then
	echo "FAIL: YepYep=No expected 'false', got '$got'"
	exit 1
fi

echo PASS
exit 0
