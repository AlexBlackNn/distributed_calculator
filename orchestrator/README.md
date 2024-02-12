# Orchestrator

docker-compose scale worker=3

Для генерации swagger документации
``` bash 
cd orchestrator
go generate ./...
```

``` bash 
go run internal/app/app.go --config=./config/local.yaml
```
``` bash 
go run main.go --config=./config/local.yaml
```

http://localhost:8080/swagger/index.html


``` bash 
curl -X POST http://localhost:8080/expression -H "Content-Type: application/json" -d '{"expression": "2*2+3"}'
``` 

``` 
export PATH=$(go env GOPATH)/bin:$PATH
cd distributed_calculator/orchestrator/cmd/orchestrator
swag init
``` 

Cannot use ginSwagger.WrapHandler(swaggerFiles.Handler)
https://github.com/swaggo/gin-swagger/issues/16#issuecomment-512813933

http://localhost:8080/swagger/index.html



## Задачи 
1. [x] 1
2. [x] 2


## Окружение развёртывания программного обеспечения - локально

#### Запуск
1. Переходим в папку с инфраструктурой и запускаем docker-compose
``` bash 
cd infra
docker-compose up
```

2. Из корня проекта накатываем миграции
```bash
go run ./cmd/migrator/postgres  --migrations-path=./migrations 
```
Примечание: в случае ошибки, подождать когда все контейнеры запустяться
```
panic: EOF
goroutine 1 [running]:
main.main()
.../sso/cmd/migrator/postgres/main_postgres.go:43 +0x29c
exit status 2
```
3. Запускаем приложение локально
```bash
go run ./main.go --config=./config/local.yaml
```

## Окружение развёртывания программного обеспечения - ДЕМО

#### Запуск
1. Переходим в папку с инфраструктурой и запускаем docker-compose 
``` bash 
cd infra
docker-compose -f docker-compose.demo.yaml up --force-recreate --build
```

2. Из корня проекта накатываем миграции
```bash
go run ./cmd/migrator/postgres  --migrations-path=./migrations 
```
Примечание: в случае ошибки, подождать когда все контейнеры запустяться
```
panic: EOF
goroutine 1 [running]:
main.main()
.../sso/cmd/migrator/postgres/main_postgres.go:43 +0x29c
exit status 2
```


