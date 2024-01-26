generate-server:
	oapi-codegen  -config ./api/quotes/server-cfg.yaml ./api/quotes/quotes.yaml

generate-models:
	oapi-codegen  -config ./api/quotes/models/models-cfg.yaml ./api/quotes/models/models.yaml

generate-all: generate-models generate-server

DOCKER_COMPOSE_FILE_PATH = ./deployments/docker-compose.yml
app-start:
	docker-compose -f $(DOCKER_COMPOSE_FILE_PATH)  up  --build -d

app-stop:
	docker-compose -f $(DOCKER_COMPOSE_FILE_PATH) stop
