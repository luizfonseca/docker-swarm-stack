deploy:
	-@docker network create --scope swarm --driver=overlay --attachable public
	-@docker network create --scope swarm --driver=overlay --attachable agent_network
	-@docker network create --scope swarm --driver=overlay --attachable monitoring
	-@docker network create --scope swarm --driver=overlay --attachable traefik-metrics

	docker stack deploy -c compose.yml -c 0-agents.yml -c 1-logging.yml -c 2-dashboard.yml -c 3-tracing.yml -c 4-registry.yml -c 5-swarm-management.yml olc

destroy:
	docker stack rm olc
	docker network rm public
	docker network rm agent_network
	docker network rm monitoring
	docker network rm traefik-metrics