#! /bin/bash
mkdir apilogs && touch apilogs/myapp.log
touch .env 
echo FULLDOMAIN=https://www.trustpilot.com/review/www.leagueoflegends.com >> .env
echo ALLOWED1=trustpilot.com >> .env
echo ALLOWED2=www.trustpilot.com >> .env
echo MAXPAGE=26 >> .env
go test -v -cover .
rm -rf ./apilogs leaguereviews.json .env
