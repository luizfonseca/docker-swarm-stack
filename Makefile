.SILENT:
.PHONY: deploy
deploy:
	@echo "Creating networks"
	@-docker network create --scope swarm --driver=overlay --attachable public 2>/dev/null ||:
	@-docker network create --scope swarm --driver=overlay --attachable agent_network 2>/dev/null ||:
	@-docker network create --scope swarm --driver=overlay --attachable monitoring 2>/dev/null ||:
	@-docker network create --scope swarm --driver=overlay --attachable traefik-metrics 2>/dev/null ||:
	@-docker network create --scope swarm --driver=overlay --attachable apps 2>/dev/null ||:
	@-docker network create --scope swarm --driver=overlay --attachable shared 2>/dev/null ||:

	./bin/deploy.sh

build-caddy:
	docker buildx build -f dockerfiles/Caddy.Dockerfile --platform linux/arm64,linux/amd64 -t luizfonseca/caddy-proxy-with-plugins:v1.0.0 ./dockerfiles/ --push

destroy:
	docker stack rm olc
	docker network rm public
	docker network rm agent_network
	docker network rm monitoring
	docker network rm traefik-metrics
	docker network rm apps
