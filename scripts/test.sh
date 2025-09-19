#!/usr/bin/env bash
set -e

################################################################################
# Usage examples:
# ./test.sh ./contrib/auth                    # uses .env.development
# ENV_PATH=.env.production ./test.sh ./contrib/auth   # uses .env.production
################################################################################

# Default env file
ENV_PATH="${ENV_PATH:-.env.test}"

echo "Using env file: $ENV_PATH"

# Load variables safely
if [ -f "$ENV_PATH" ]; then
  while IFS='=' read -r key value; do
    # skip comments and empty lines
    [[ "$key" =~ ^#.*$ || -z "$key" ]] && continue
    # remove surrounding quotes from value if any
    value="${value%\"}"
    value="${value#\"}"
    export "$key=$value"
  done < "$ENV_PATH"
else
  echo "Env file not found: $ENV_PATH"
  exit 1
fi

# Run go test with passed args
go test -v "$@"
