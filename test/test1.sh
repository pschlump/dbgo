#!/bin/bash

go build

export YepYep=Yes
rv=$(./test)
# echo "rv=$rv"
if [ "$rv" == "false" ] ; then
	echo FAIL: "failed to return true when it should"
	exit 1
fi

export YepYep=No
rv=$(./test)
# echo "rv=$rv"
if [ "$rv" == "true" ] ; then
	echo FAIL: "failed to return false when it should"
	exit 1
fi

echo PASS

exit 0

