export GOPATH=`pwd`
cd src/tourney
rm tourney
go build
./tourney
cd ../..