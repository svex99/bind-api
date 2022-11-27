#!/bin/sh

# Abort on any error (including if wait-for-it fails).
set -e

# Wait for the backend to be up, if we know where it is.
if [ -n "$BIND_HOST" ]; then
  ./wait-for-it.sh "$BIND_HOST:${BIND_PORT:-53}"
fi

# Run the main container command.
exec "$@"
