generate-server:
	oapi-codegen  -config ./api/quotes/server-cfg.yaml ./api/quotes/quotes.yaml

generate-models:
	oapi-codegen  -config ./api/quotes/models/models-cfg.yaml ./api/quotes/models/models.yaml

generate-all: generate-models generate-server

app-start:
	docker-compose up  --build -d

app-stop:
	docker-compose stop
