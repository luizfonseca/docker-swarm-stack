#! /bin/bash

echo "$(tput setaf 1)IMPORTANT: $(tput sgr0)"
echo " "
echo "$(tput setaf 1)In order for container metrics to be scraped properly $(tput sgr0)"
echo "$(tput setaf 1)we need to create a "socat" container that creates a TCP host:port $(tput sgr0)"
echo "$(tput setaf 1)bind for the /var/run/docker.sock socket. $(tput sgr0)"

echo " "
echo " "
echo " "
read -n 1 -p "Confirm? (y/n):" "mainmenuinput"

if [ "$mainmenuinput" = "y" ]; then
  echo " "
  echo "$(tput setaf 4)>> Creating socat container$(tput sgr0)"

  # mkdir -p /var/docker-socat
  docker service rm docker-metrics
  docker service create \
    -d \
    --name docker-metrics \
    --mode global \
    --mount type=bind,source=/var/run/docker.sock,destination=/var/run/docker.sock \
    --network host \
    --publish 2376:2376 \
    alpine/socat \
    tcp-listen:2376,fork,reuseaddr unix-connect:/var/run/docker.sock

  echo "$(tput setaf 4)>> Socat container created $(tput sgr0)"
else
  echo "$(tput setaf 1)Aborting$(tput sgr0)"
  exit 1
fi

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