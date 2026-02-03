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

