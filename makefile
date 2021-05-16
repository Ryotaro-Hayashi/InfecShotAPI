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
	docker exec -it infecshot-mysql sh

.PHONY: infecshot-api
api:
	docker exec -it infecshot-api sh

server:
	docker exec infecshot-api go run cmd/main.go

stop:
	docker exec infecshot-api pkill -e go
	docker exec infecshot-api pkill -e main
