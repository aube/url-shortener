#!/bin/bash

# if [ ]

if [ -z $1 ]; then
    BRANCH=$(git branch | sed -n -e 's/^\* \(.*\)/\1/p')
    ITERATION=${BRANCH: -1}
else
    ITERATION=$1
fi

./shortenertest \
 -test.v \
 -test.run=^TestIteration$ITERATION$ \
 -binary-path=cmd/shortener/shortener \
 -server-port=8080 \
 -file-storage-path=_hashes/hashes_list.json \
 -source-path=. \
 -database-dsn=123


echo "done iteration:" $ITERATION