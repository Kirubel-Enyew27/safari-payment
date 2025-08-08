migrate-down:
	- migrate -database postgres://user:password@localhost:5432/db?sslmode=disable -path internal/query/schemas -verbose down
migrate-up:
	- migrate -database postgres://user:password@localhost:5432/db?sslmode=disable -path internal/query/schemas -verbose up
migrate-create:
	- migrate create -ext sql -dir internal/query/schemas -tz "UTC" $(args)
sqlc:
	- sqlc generate -f ./config/sqlc.yaml
