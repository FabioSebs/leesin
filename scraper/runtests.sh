#! /bin/bash
mkdir apilogs
go test -v -cover .
rm -rf ./apilogs leaguereviews.json
