#!/bin/bash
rm tourney 2>/dev/null
cd src
go build || exit 1
cd ..
mv src/src tourney || exit 1
./tourney 
