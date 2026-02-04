<p style="display:none">TODO: Create a unique logo for the application</p>
<div align="center">
  <img height="150" alt="market-management-app-logo" src="https://"/>
  <h1>🛜 Market Management Application Server API</h1>
  <p>Modular monolithic backend application for market management app</p>

![Static Badge](https://img.shields.io/badge/Go-blue?style=plastic&logo=go&logoColor=ffffff)
![Static Badge](https://img.shields.io/badge/Docker-blue?style=plastic&logo=docker&logoColor=ffffff)
![Static Badge](https://img.shields.io/badge/Prometheus-E6522C?style=plastic&logo=prometheus&logoColor=white)
![Static Badge](https://img.shields.io/badge/OpenTelemetry-ffffff?style=plastic&logo=opentelemetry&logoColor=000000)
![Static Badge](https://img.shields.io/badge/Redis-FF4438?style=plastic&logo=redis&logoColor=white)
![Static Badge](https://img.shields.io/badge/PostreSQL-4169E1?style=plastic&logo=postgresql&logoColor=white)
![Static Badge](https://img.shields.io/badge/MongoDB-white?style=plastic&logo=mongodb&logoColor=#47A248&labelColor=#47A248)

![GitHub last commit](https://img.shields.io/github/last-commit/ilkerciblak/dukkanim-app)
![GitHub last commit](https://img.shields.io/github/contributors/ilkerciblak/dukkanim-app)

</div>

## Project Introduction and Architectural Decisions


## Pre-requisities and Project Setup

### Pre-requisities
- Project's environment only requires a pre-ready `Docker` installation on local computer.

#### Configuring Environment Variables
Various project platform utilites depends on the environment variables that should be defined in `.env` file in the projects' directory. For an example, data persisting platform requires `postgres connection string` or in order to cache with `redis` implementation `redis host and port` should be declared before project run-time.

- Create `project-folder/backend/.env`, copy following code and declare variables
```yaml
# APP
APP_MODE=DEV
APP_PORT=8000


# DATABASE
POSTGRES_USER=admin
POSTGRES_DB=db
POSTGRES_PASSWORD=admin
DB_HOST=postgres-db:5432
CONN_STR=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${DB_HOST}/postgres?sslmode=disable

# GOOSE
GOOSE_DRIVER=postgres
GOOSE_DBSTRING=${CONN_STR}
GOOSE_MIGRATION_DIR=migrations/*/*


# Rate Limiter
RateLimiter_TimeFrame_Seconds=60
RateLimiter_Request_Per_TimeFrame=100


# Observability
LOG_LEVEL=ERROR

# REDIS
REDIS_HOST=go-redis
REDIS_PORT=6379
REDIS_ADDR=${REDIS_HOST}:${REDIS_PORT}

# MongoDB
MONGO_INITDB_ROOT_USERNAME=root
MONGO_INITDB_ROOT_PASSWORD=example
MONGO_URL=mongodb://${MONGO_INITDB_ROOT_USERNAME}:${MONGO_INITDB_ROOT_PASSWORD}@go-mongo?retryWrites=true

# Mongo-Express
ME_CONFIG_MONGODB_URL=mongodb://root:example@go-mongo:27017/
ME_CONFIG_BASICAUTH_ENABLED=true
ME_CONFIG_BASICAUTH_USERNAME=mongoexpressuser
ME_CONFIG_BASICAUTH_PASSWORD=mongoexpresspass
```



1. Clone the project repository in local computer
``` bash
$ git clone https://github.com/ilkerciblak/dukkanim-app
```
> [!IMPORTANT]
> Repository clonning is not required also redundant if you done it already in [Project Frontend Part]()
<br/>

2. Change directory to this directory
```bash
$ cd project_folder/backend
```

3. Build and run backend container using docker compose file
```bash
$ docker compose -f ../dev.docker-compose.yml up -d
```

## Repository Structure

```
.
├── backend
│   ├── bin
│   │   └── app
│   ├── cmd
│   │   └── api
│   │       └── main.go 
│   ├── internal
│   │   ├── api
│   │   │   ├── http.go
│   │   │   └── middleware/
│   │   ├── features
│   │   │   ├── product/
│   │   │   │   ├── domain
│   │   │   │   │   ├── product_categories.go
│   │   │   │   │   ├── product_unit-types.go
│   │   │   │   │   └── product.go
│   │   │   │   ├── internal
│   │   │   │   │   ├── handler.go
│   │   │   │   │   ├── repository._go
│   │   │   │   │   └── service._go
│   │   │   │   └── router.go
│   │   │   └── ...other-features/
│   │   └── platform
│   │       ├── caching/
│   │       ├── config/
│   │       ├── database/
│   │       ├── http_response/
│   │       ├── observability/
│   │       ├── problem/
│   │       ├── rate_limiting/
│   │       └── timestamp/
│   ├── pkg/
│   │   └── viladition/
│   │       ├── internal
│   │       ├── viladitor_test.go
│   │       └── viladitor.go
│   ├── migrations/
│   │   
│   ├── go.mod
│   ├── go.sum
│   ├── Makefile
│   ├── prometheus.yml
│   ├── backend.Dockerfile
│   └── README.md
```

### Repository Architecture Overview

#### Platform Directory

#### Features Directory

#### Internal/api Directory

## API Documentations

## Repository Decision Records

## Contact





