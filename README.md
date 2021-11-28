# golang-rest-api
It's a simple REST API application written in Golang using GORM and PostgreSQL. To run you need [docker](https://www.docker.com/) to be installed  
## How to run
Build and run
```
build-run
```
It automatically migrates up as first launched so there is no need to make a command for it
```
make migration-up
```
You can also migration-down
```
make migration-down
```
Stop all containers
```
make stop
```
## In addition
Run tests
```
make test
```
golangci-lint
```
make lint
```
