# Chat
A service to help people communication in VDAT ecosystem

## Quickstart
> clone dependencies
```shell script
git submodule update --init
```


> run with Docker Compose
```shell script
docker-compose up
```

> set up .env
```shell script
touch .env
echo MINIO_END_PORT=localhost:9000 >> .env
echo MINIO_KEY=minio >> .env
echo MINIO_ACCESSES=minio123 >> .env
echo DATABASE_URL=postgres://postgres:postgres@localhost:5432/dchat >> .env
```

> install library and dependencies
```shell script
 go mod tidy
```

> run server
```shell script
 go run ./cmd/chatserver
```

> get token from keycloak
```shell script
 go run ./cmd/chatcli
```

> run website
```shell script
 cd website
 npm install (only first time to install library)
 npm start
```

## Environments
* Production: https://vdat-mcsvc-chat.vdatlab.com
* Staging: https://vdat-mcsvc-chat-staging.vdatlab.com/


## Swagger
* Production: https://vdat-mcsvc-chat.vdatlab.com/swagger/index.html#/
* Staging: https://vdat-mcsvc-chat-staging.vdatlab.com/swagger/index.html#/
* Local: http://localhost:5000/swagger/index.html#/



see full list [here](https://gitlab.com/vdat/mcsvc/chat/-/environments).

## Plans
### Version 0.1
Basic functionality, including:
1. Basic UI interface
2. View messages and send message
3. Searching people

### Version 0.2
Basic functionality, including:
1. Basic UI interface
2. View messages and send message
3. Searching people and add user to group
4. article for Q&A (User can write or comment, reaction)

## ERD diagram
![](docs/Domain%20Objects.png)

## Schedule
![](docs/schedule.png)

###### to update schedule modify `./docs/schedule.puml` and save the result in `./docs/schedule.png`

## Architecture
![](docs/architecture.png)

`Chat Service` will dependent on some other services:
1. `Identity Provider` for search feature
2. `Authz Server` for `access token` verification

## Development

### Monitoring tools
use following command to run Prometheus, Grafana to monitor performance
```shell script
docker-compose -f dev/docker-compose.yml up
```
* access Grafana's Dashboards at http://localhost:3000/dashboards

### Project layout guideline

`/cmd` contains binary package
`/pkg` contains shared modules

see https://github.com/golang-standards/project-layout
