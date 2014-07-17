cd src/bitboard
go build
go install
cd ../tourney
rm tourney
go build
./tourney
cd ../..