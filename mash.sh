#!/bin/bash

export DATABASE_URL="postgres://admin:adminpass@127.0.0.1:5432/game?sslmode=disable"
export RELEASE_MODE=dev

set -e

BASE_LOCATION=$(pwd)

generate_sqlc() {
	echo "Generating..."
	cd ./backend/
	sqlc generate
	cd ..
	echo "OK"	
}

COMMAND=$1

if [ -z "COMMAND" ]; then
    echo "i need a command"
fi

case $COMMAND in
	run)
		generate_sqlc
		echo "Starting Core"
		cd ./backend/
		go run ./cmd/api/.
		cd ..
		echo "The End"
		;;
	sqlc)
		generate_sqlc
		;;
	run-front)
		cd ./frontend/
		npm run dev
		cd ..
		;;
	*)
		echo "Invalid command"
		;;
esac
