#! /bin/bash
echo "$(tput setaf 4)>> Loading environment variables from .env file$(tput sgr0)"
eval "$(
  cat .env | awk '!/^\s*#/' | awk '!/^\s*$/' | while IFS='' read -r line; do
    key=$(echo "$line" | cut -d '=' -f 1)
    value=$(echo "$line" | cut -d '=' -f 2-)
    echo "export $key=\"$value\""
  done
)"
echo "$(tput setaf 4)>> Deploying services$(tput sgr0)"


tput setaf 6 && docker stack deploy -c main.compose.yml -c 0-agents.compose.yml -c 1-logging.compose.yml -c 2-dashboard.compose.yml -c 3-tracing.compose.yml -c 4-registry.compose.yml -c 5-swarm-management.compose.yml olc


echo " "
echo "$(tput setaf 2)The following services URLs were created by default:$(tput sgr0)"
echo " - $(tput setaf 3)https://portainer$DOMAIN_SUFFIX.$DOMAIN_NAME$(tput sgr0) (Stack Management)"
echo " - $(tput setaf 3)https://grafana$DOMAIN_SUFFIX.$DOMAIN_NAME$(tput sgr0) (Monitoring, Logs, Traces)"
echo " - $(tput setaf 3)https://registry$DOMAIN_SUFFIX.$DOMAIN_NAME$(tput sgr0) (Private Docker Registry)"
echo " - $(tput setaf 3)https://gh-oauth$DOMAIN_SUFFIX.$DOMAIN_NAME$(tput sgr0) (Github Oauth Proxy)"
echo " "