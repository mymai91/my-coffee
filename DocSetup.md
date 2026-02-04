## Step 1: Project Structure

```
# Create project directory
mkdir my-new-api
cd my-new-api

# Initialize Go module
go mod init github.com/yourusername/my-new-api

```

## Step 2: Install Dependencies (First Time Only)

### Core Dependencies:

```
# Database & ORM
go get gorm.io/gorm
go get gorm.io/driver/postgres

# Configuration
go get github.com/joho/godotenv

# Validation
go get github.com/go-playground/validator/v10

# Migration dependencies (for cmd/migrate)
go get github.com/golang-migrate/migrate/v4
go get github.com/golang-migrate/migrate/v4/database/postgres
go get github.com/golang-migrate/migrate/v4/source/file
go get github.com/lib/pq
```

### Development Tools (Install Globally):

```
# Migration CLI tool
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

## Step 3: Environment variables
### Create .env

```
# Database
DB_HOST=localhost
DB_PORT=5454
DB_USER=jany
DB_PASSWORD=mypassword
DB_NAME=coffee_dev
DB_SSL_MODE=disable
```

## Step 4: Create Config & DB setup
### Open config/config.go

```
```

## Docker
### Open docker-compose.yml

```
version: '3.8'
services:
  dev-db:
    image: postgres:15
    ports:
      - 5454:5432
    environment:
      - POSTGRES_USER=jany
      - POSTGRES_PASSWORD=mypassword
      - POSTGRES_DB=coffee_dev
#     volumes:
#       - dev-db-data:/var/lib/postgresql/data
# volumes:
#   dev-db-data: 
#     driver: local
```


# Part 2: Migration Management

### 1. Start PostgreSQL
make docker-db-dev-up

### 2. Create a new migration
make migrate-create name=create_orders_table

### 3. Run migrations
make migrate-up

### 4. Check status
make migrate-status