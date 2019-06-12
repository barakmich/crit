#!/bin/bash

set -xv

rm -rf /tmp/foo
export CRIT_REVIEW_REPO=/tmp/foo
./crit init
#./crit start example --base-url $PWD/examplerepo --review-branch simple
./crit start example --base-url $PWD --review-branch master --base-branch base
./crit list 
./crit open example
