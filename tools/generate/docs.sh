#!/bin/sh
DOCS="../../docs"
MP="$DOCS/man"
MD="$DOCS/markdown"

mkdir -p $MP $MD

go run docs.go $DOCS

tar -czf $MP/fox.tar.gz $MP/*.1
rm $MP/*.1

exit 0
