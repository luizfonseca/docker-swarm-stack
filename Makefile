.SILENT:
.PHONY: deploy
deploy:
	@-docker network create --scope swarm --driver=overlay --attachable public 2>/dev/null ||:
	@-docker network create --scope swarm --driver=overlay --attachable agent_network 2>/dev/null ||:
	@-docker network create --scope swarm --driver=overlay --attachable monitoring 2>/dev/null ||:
	@-docker network create --scope swarm --driver=overlay --attachable traefik-metrics 2>/dev/null ||:

	./bin/deploy.sh

destroy:
	docker stack rm olc
	docker network rm public
	docker network rm agent_network
	docker network rm monitoring
	docker network rm traefik-metrics