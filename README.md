# Go Clean Template

## Description

This is Challenge Backend 1 template.

## Architecture

![Clean Architecture](architecture.png)

1. External system perform request (HTTP, gRPC, Messaging, etc)
2. The Delivery creates various Model from request data
3. The Delivery calls Use Case, and execute it using Model data
4. The Use Case create Entity data for the business logic
5. The Use Case calls Repository, and execute it using Entity data
6. The Repository use Entity data to perform database operation
7. The Repository perform database operation to the database
8. The Use Case create various Model for Gateway or from Entity data
9. The Use Case calls Gateway, and execute it using Model data
10. The Gateway using Model data to construct request to external system
11. The Gateway perform request to external system (HTTP, gRPC, Messaging, etc)

## Tech Stack

- Golang : https://github.com/golang/go
- PostgreSQL (Database) : https://github.com/postgres/postgres
- Apache Kafka : https://github.com/apache/kafka

## Framework & Library

- GoFiber (HTTP Framework) : https://github.com/gofiber/fiber
- GORM (ORM) : https://github.com/go-gorm/gorm
- Viper (Configuration) : https://github.com/spf13/viper
- Golang Migrate (Database Migration) : https://github.com/golang-migrate/migrate
- Go Playground Validator (Validation) : https://github.com/go-playground/validator
- Zap (Logger) : https://github.com/uber-go/zap
- Sarama (Kafka Client) : https://github.com/IBM/sarama

## Configuration

### Files

- `config.json`: Application defaults.
- `.env`: Secrets and local overrides (if needed).

### Environment Variables

Make sure to rename `env.example` and `config.json.example` to `.env` and `config.json`, respectively.

You can override any config in `config.json` using environment variables. The key mapping uses `_` instead of `.`.
Example `.env`:

```env
API_KEY=your-secret-key-123
```

Ensure you create a `.env` file before running the application. Use `.env.example` as a template if available.

## API Spec

All API Spec is in `docs` folder.

## Docker Setup

### Up containers (db, broker, and zookeeper)

```shell
make docker-up
```

### Down containers (db, broker, and zookeeper)

```shell
make docker-down
```

With volume:

```shell
make docker-down-v
```

## Database Migration

All database migration is in `db/migrations` folder.

### Create Migration

```shell
make migrate-create name=create_table_xxx
```

### Run Migration

```shell
make migrate-up
```

### Rollback Migration

```shell
make migrate-down
```

### Reset Migration

```shell
make migrate-reset
```

### Generate API Spec

```shell
make swag
```

## Run Application

### Run unit test

```bash
make test
```

### Run web server

```bash
make run
```

### Run worker

```bash
make run-worker
```

### Hot Reload

To run the application with hot reload, use [Air](https://github.com/air-verse/air).

```bash
make dev
```

## Contact Me

**OPEN FOR CONTRIBUTIONS.**

This repository is an improved version from [khannedy's Clean Architecture](https://github.com/khannedy/challenge-backend-1). If you have any questions or suggestions, please feel free to contact me.

- Email : achievafuturagemilang@gmail.com
- LinkedIn : https://www.linkedin.com/in/achieva-futura-gemilang
- GitHub : https://github.com/achievagemilang

## Reference

- [khannedy's Clean Architecture](https://github.com/khannedy/challenge-backend-1)
