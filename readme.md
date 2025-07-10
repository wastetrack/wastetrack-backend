# Generating Migrations
migrate create -ext sql -dir database/migrations -seq '<migration_name>'

# Apply Migrations
migrate -database "postgres://user:password@localhost:5432/your_db?sslmode=disable" -path database/migrations up 

# Rollback
1. to undo migrations
   ```
   migrate -database "postgres://user:password@localhost:5432/your_db?sslmode=disable" -path database/migrations down 1
   ```

2. to reset migrations
   ```
   migrate -database "postgres://user:password@localhost:5432/your_db?sslmode=disable" -path database/migrations down
   ```

3. Force state
   ```
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
go run .\database\cmd\seeder\main.go -db "postgres://postgres:Ryokeren123@localhost:5432/wastetrack?sslmode=disable"
```
### Get list of commands
```bash
go run ./database/cmd/seeder/main.go -help
```
### Clear seeders
```bash
go run .\database\cmd\seeder\main.go -db "postgres://postgres:Ryokeren123@localhost:5432/wastetrack?sslmode=disable" -clear
```

# Run web server

```bash
go run cmd/web/main.go
```

