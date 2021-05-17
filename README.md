# infecshot API

### Installation
1. copy `.env.example` file, and make `.env` file
2. $`docker network create loki`
3. $`make local`
4. $`make server`

### Stop
$`make stop`

### Test
$`docker-compose -f docker-compose.yml -f docker-compose.local.yml run --rm api go test ./...`
