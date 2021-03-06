#!/bin/sh

set -exu

BIN="percona-agent"
CWD="$PWD"

./agent-build/agent-build

cd ../bin/$BIN-installer/
go build

cd $CWD

[ -f $BIN.tar.gz ] && rm -f $BIN.tar.gz
if [ -d $BIN ]; then
   rm -rf $BIN/*
fi
mkdir -p "$BIN/bin" "$BIN/init.d"

cp ../install/install.sh $BIN/install
cp ../COPYING ../README.md ../Changelog.md ../Authors $BIN/
cp ../bin/$BIN/$BIN ../bin/$BIN-installer/$BIN-installer $BIN/bin
cp ../install/$BIN $BIN/init.d

tar cvfz $BIN.tar.gz $BIN/
