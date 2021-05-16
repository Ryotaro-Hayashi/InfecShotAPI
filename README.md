# infecshot API

### Installation
1. copy `.env.example` file, and make `.env` file
2. $`make local`
3. $`make server`

### Stop
$`make stop`

### Test
$`docker-compose -f docker-compose.yml -f docker-compose.local.yml run --rm api go test ./...`
