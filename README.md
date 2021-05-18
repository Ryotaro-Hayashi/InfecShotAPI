# infecshot API

### Installation
1. $`cp .env.example .env`
2. $`docker network create loki`
3. $`make local`
4. $`make server`

### Stop
$`make stop`

### Test
$`docker-compose -f docker-compose.yml -f docker-compose.local.yml run --rm api go test ./...`
