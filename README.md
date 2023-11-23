# OLC (One-Line _Self-Hosting_ Command)


This repository is a collection of scripts and configuration files to setup a self-hosting stack on a single server (or multiple) using `docker swarm` as the main orchestrator with a simple one-line command.

The stack is composed of the following services:
- [Traefik](https://traefik.io/) as a reverse proxy
- [Portainer](https://www.portainer.io/) as a docker management UI (and service deployment)
- [Docker Registry](https://hub.docker.com/_/registry) to host your **own** docker images and trigger deployments through Github Actions etc.
- Grafana to monitor the server's resources and logging
- Loki to collect logs from all the services
- Vector to act as a log sink for Loki and send logs to Grafana (or other places you define)


## Requirements

#### A server with a public IP address. 

Check out [Hetzner](https://www.hetzner.com/), [Contabo](https://contabo.com/en/), [Vultr](https://www.vultr.com/), [DigitalOcean](https://www.digitalocean.com/), [Linode](https://www.linode.com/) for servers within 5 EUR/USD.

#### A domain name
A domain name and the correct DNS records pointing to the server. This is particular to your DNS provider.  Often you can find these settings in the DNS management panel.

* In most cases, you will need to create an `A` record pointing to the server's IP address. [More about it here.](https://www.cloudflare.com/learning/dns/dns-records/dns-a-record/)

* In most cases, you will also need to create a `CNAME` record pointing to the `A` record. [More about it here.](https://www.cloudflare.com/learning/dns/dns-records/dns-cname-record/)

#### A DNS provider that supports ACME (Let's Encrypt)

[ACME](https://en.wikipedia.org/wiki/Automated_Certificate_Management_Environment) (Let's Encrypt)

#### Docker and docker-compose installed on the server 

We will use `docker swarm` to manage the services together. I recommend using the **latest stable** version of docker and docker-compose.

Installation instructions [can be found here for Docker Engine](https://docs.docker.com/engine/install/). 
If you are using `Ubuntu`, you can use the following link: https://docs.docker.com/engine/install/ubuntu/

## Usage

1. Clone this repository on the server (or download)
2. Make sure you have a swarm initialized on the server. If not, run the following command:

```bash
docker swarm init
```

3. Run the following command to start the stack:

```bash
./oshc.sh --domain <YOUR_DOMAIN> --email <YOUR_EMAIL> \ 
    --registry-username <REGISTRY_USERNAME> --registry-password <REGISTRY_PASSWORD> \


```