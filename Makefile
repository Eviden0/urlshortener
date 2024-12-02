lanch_postgres:
	docker run --name postgres-url \
	-e POSTGRES_USER=lang \
	-e POSTGRES_PASSWORD=password \
	-e POSTGRES_DB=urldb \
	-p 5432:5432 \
	-d postgres

lanch_redis:
	docker run --name redis \
	-p 6379:6379 \
	-d redis

migrate_up:
	migrate -path="./db/migration" -database="postgres://lang:password@localhost:5432/urldb?sslmode=disable" up
	
migrate_down:
	migrate -path="./db/migration" -database="postgres://lang:password@localhost:5432/urldb?sslmode=disable" drop -f