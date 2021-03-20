up:
	docker-compose up -d

down:
	docker-compose down -v

mysql:
	docker exec -it mysql sh

.PHONY: api
api:
	docker exec -it api sh

server:
	go run cmd/main.go
