#!/bin/sh
set -e
goose -dir ./migrations up ${GOOSE_DRIVER} "${GOOSE_DBSTRING}" up
exec ./main