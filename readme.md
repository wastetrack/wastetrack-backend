# Generating Migrations
migrate create -ext sql -dir database/migrations -seq '<migration_name>'

# Apply Migrations
migrate -database "postgres://user:password@localhost:5432/your_db?sslmode=disable" -path database/migrations up 

# Rollback
1. to undo migrations
   ```bash
   migrate -database "postgres://user:password@localhost:5432/your_db?sslmode=disable" -path database/migrations down 1
   ```

2. to reset migrations
   ```bash
   migrate -database "postgres://user:password@localhost:5432/your_db?sslmode=disable" -path database/migrations down
   ```

3. Force state
   ```bash
   migrate -database "postgres://user:password@localhost:5432/your_db?sslmode=disable" -path database/migrations force 0
   migrate -database "postgres://user:password@localhost:5432/your_db?sslmode=disable" -path database/migrations force 1
   ```

# First time running
## Getting packages
```bash
go get ./...
go go mod tidy
```
## Seeder
### Apply seeders
```bash
go run .\database\cmd\seeder\main.go -db "postgres://user:password@localhost:5432/your_db?sslmode=disable"
```
### Get list of commands
```bash
go run ./database/cmd/seeder/main.go -help
```
### Clear seeders
```bash
go run .\database\cmd\seeder\main.go -db "postgres://user:password@localhost:5432/your_db?sslmode=disable" -clear
```

# Run web server

```bash
go run cmd/web/main.go
```

# Run using Makefile

1. Copy Makefile.example ke Makefile.

   ```bash
   cp Makefile.example Makefile
   ```

2. Edit `DB_URL_LOCAL` in Makefile with your configuration settings.

3. Read this documentation to run what you want. :D

   ```bash
   make create-migration     # make new migration
   make migrate-up           # apply migration
   make migrate-down         # rollback 1 step
   make migrate-reset        # reset all
   make migrate-force        # force specific version

   make deps                 # run dependency
   make seed                 # generate data dummy
   make seed-clear           # delete data dummy
   make seed-help            # see seeder command

   make run-local            # run the local server

   make run-vps              # run the server on VPS (background)
   make stop-vps             # stop the VPS server
   ```