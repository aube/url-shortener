#!/bin/bash

if [ -z $1 ]; then
    ITERATION=$(git branch | sed -n -e 's/^\* \(.*\)/\1/p' | tr -d -c 0-9)
else
    ITERATION=$1
fi

cd cmd/shortener

go build -o shortener *.go

cd -

./shortenertest \
 -test.v \
 -test.run=^TestIteration$ITERATION$ \
 -binary-path=cmd/shortener/shortener \
 -server-port=8080 \
 -file-storage-path=_hashes/hashes_list.json \
 -source-path=. \
 -database-dsn=123


echo "done iteration:" $ITERATION
