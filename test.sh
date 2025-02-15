#!/bin/bash

BRANCH=$(git branch | sed -n -e 's/^\* \(.*\)/\1/p')
ITERATION=${BRANCH: -1}

./shortenertest -test.v -test.run=^TestIteration$ITERATION$ -binary-path=cmd/shortener/shortener -server-port=8080


