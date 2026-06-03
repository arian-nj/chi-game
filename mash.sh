#!/bin/bash

export DATABASE_URL="postgres://admin:adminpass@127.0.0.1:5432/game?sslmode=disable"
export RELEASE_MODE=dev
export JWT_SECRET="jjlkds5r6789uipokkjhj3t7y8u9psty543sghh0jlkkd5r6t7y8u"

set -e

BASE_LOCATION=$(pwd)

generate_sqlc() {
	echo "Generating..."
	cd ./backend/
	sqlc generate
	cd ..
	echo "OK"
}

run_backend() {
	generate_sqlc
	echo "Starting backend on :8383"
	cd ./backend/
	go run ./cmd/api/
	cd ..
}

COMMAND=$1

if [ -z "$COMMAND" ]; then
	echo "Usage: ./mash.sh {run-back|run-front|sqlc|run}"
	exit 1
fi

case $COMMAND in
	run|run-back)
		run_backend
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
