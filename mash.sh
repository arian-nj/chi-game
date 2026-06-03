#!/bin/bash

export DATABASE_URL="postgres://admin:adminpass@127.0.0.1:5432/game?sslmode=disable"
export RELEASE_MODE=dev
export JWT_SECRET="jjlkds5r6789uipokkjhj3t7y8u9psty543sghh0jlkkd5r6t7y8u"
export ADMIN_SECRET="b12v44DSacdxserte5r467t78youijDSFHEYTWiuyuehdhyhqwdd"
export VITE_API_BASE_URL="https://api.chigame.site"
export VITE_SITE_URL="https://chigame.site"

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

run_tests() {
	generate_sqlc
	echo "Running backend tests..."
	cd ./backend/
	mapfile -t test_pkgs < <(go list -f '{{if len .TestGoFiles}}{{.ImportPath}}{{end}}' ./...)
	if [ "${#test_pkgs[@]}" -eq 0 ]; then
		echo "No packages with tests found"
		exit 1
	fi
	go test -count=1 "${test_pkgs[@]}"
	cd ..
	echo "OK"
}

COMMAND=$1

if [ -z "$COMMAND" ]; then
	echo "Usage: ./mash.sh {run-back|run-front|sqlc|test|run}"
	exit 1
fi

case $COMMAND in
	run|run-back)
		run_backend
		;;
	
	sqlc)
		generate_sqlc
		;;
	
	test)
		run_tests
		;;
	
	run-front)
		cd ./frontend/
		npm run dev
		cd ..
		;;

	deploy-back)
		cd ./backend/
		GOOS=linux GOARCH=amd64 go build -o main ./cmd/api/.
		cd ..
		sftp apichi:/home/deploy/chigame/ <<EOF
put ./backend/main
put ./release.sh
EOF

		;;

	deploy-front)
		cd ./frontend/
		npm run build
		npx wrangler pages deploy ./dist
		;;

	serve)
		cd ./frontend/
		npm run build
		npx serve ./dist
		;;
	*)
		echo "Invalid command"
		;;
esac