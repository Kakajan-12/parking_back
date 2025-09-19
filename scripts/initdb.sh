#!/usr/bin/env bash
set -e

# Default env file
ENV_PATH="${ENV_PATH:-.env.development}"

echo "Using env file: $ENV_PATH"

# # Load variables safely
# if [ -f "$ENV_PATH" ]; then
#   while IFS='=' read -r key value; do
#     # skip comments and empty lines
#     [[ "$key" =~ ^#.*$ || -z "$key" ]] && continue
#     # remove surrounding quotes from value if any
#     value="${value%\"}"
#     value="${value#\"}"
#     export "$key=$value"
#   done < "$ENV_PATH"
# else
#   echo "Env file not found: $ENV_PATH"
#   exit 1
# fi

# Run app
ENV_PATH=$ENV_PATH go run initdb.go