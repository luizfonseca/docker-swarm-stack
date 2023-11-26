#! /bin/bash

echo ">> Loading environment variables from .env file"
eval "$(
  cat .env | awk '!/^\s*#/' | awk '!/^\s*$/' | while IFS='' read -r line; do
    key=$(echo "$line" | cut -d '=' -f 1)
    value=$(echo "$line" | cut -d '=' -f 2-)
    echo "export $key=\"$value\""
  done
)"
echo "   DONE"
echo ">> Deploying services"


docker stack deploy -c compose.yml -c 0-agents.yml -c 1-logging.yml -c 2-dashboard.yml -c 3-tracing.yml -c 4-registry.yml -c 5-swarm-management.yml olc
