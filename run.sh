#!/usr/sh

set -e

git clone --depth 1 https://github.com/andrielfn/pg-ulid.git
cd pg-ulid
make install

compareids all --host postgres
