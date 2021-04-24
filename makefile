.PHONY: local
local:
	docker-compose -f docker-compose.yml -f docker-compose.local.yml up -d

.PHONY: prod
prod:
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

.PHONY: down
down:
	docker-compose down -v --remove-orphans

mysql:
	docker exec -it mysql sh

.PHONY: infecshot-api
api:
	docker exec -it infecshot-api sh

.PHONY: server
server:
	go run cmd/main.go

.PHONY: server&
server&:
	go run cmd/main.go &
