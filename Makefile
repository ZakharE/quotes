start_app:
	docker-compose up . -d
	goose migrate up

generate-server:
	oapi-codegen  -config ./api/quotes/server-cfg.yaml ./api/quotes/quotes.yaml

generate-models:
	oapi-codegen  -config ./api/quotes/models/models-cfg.yaml ./api/quotes/models/models.yaml

generate-all: generate-models generate-server
