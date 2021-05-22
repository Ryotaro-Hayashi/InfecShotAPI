# infecshot API

### client
https://ryotaro-hayashi.github.io/InfecShot/

### Structure
![infecshot-api-overall](https://user-images.githubusercontent.com/53222150/119227638-de35b300-bb49-11eb-8785-67d446d9f280.png)

### Installation
1. $`git clone https://github.com/Ryotaro-Hayashi/InfecShotAPI.git`
2. $`cd InfecShotAPI`
3. $`cp .env.example .env`
4. $`docker network create loki`

### Start
1. $`make local`
2. $`make server`

### Access
* server: http://localhost:8080/
* swagger: http://localhost:3030/
* grafana: http://localhost:3000/
* phpmyadmin: http://localhost:4000/

### Stop
$`make stop`

### Test
1. $`make local`
2. $`make test`

### Down
$`make down`
