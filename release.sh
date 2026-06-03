#!/bin/bash

set -e

# --- Configurations --- 
export DATABASE_URL="postgres://postgres:mn5uyiut4srdtfuhisd3tf3@127.0.0.1:5432/game?sslmode=disable"
export RELEASE_MODE=release
export JWT_SECRET="kjxdser567t8910002SD0opllsSwew232oksDkjahjniaolgfdre56779i1"

export ADMIN_SECRET="b12v44DSacdxserte5r467t78youijDSFHEYTWiuyuehdhyhqwdd"



COMMAND=$1

if [ -z "COMMAND" ]; then
	echo "i want more commands"
fi

case $COMMAND in
	run)
		echo "Starting Core Bot"
		./main
		echo "The End"
		;;

	*)
		echo "command '$COMMAND' is Unknown"
esac

exit 0


