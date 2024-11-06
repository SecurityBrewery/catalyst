#!/bin/bash

# Set the flags
FLAGS=""
if [ -n "$CATALYST_FLAGS" ]; then
    FLAGS="$CATALYST_FLAGS"
fi

# Set the app url
APP_URL=""
if [ -n "$CATALYST_APP_URL" ]; then
    APP_URL="$CATALYST_APP_URL"
fi

/usr/local/bin/catalyst serve --http 0.0.0.0:8080 --flags "$FLAGS" --app-url "$APP_URL"