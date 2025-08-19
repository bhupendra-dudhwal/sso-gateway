APP_NAME=tyroo-central

.PHONY: build run stop logs db-psql redis-cli clean \
        deps deps-up deps-down \
        postgres-up keydb-up \
        postgres-down keydb-down

build:
	docker compose build

run:
	docker compose up -d

stop:
	docker compose down

logs:
	docker compose logs -f $(APP_NAME)

db-psql:
	docker exec -it postgres psql -U admin -d tyroo-central

redis-cli:
	docker exec -it keydb keydb-cli -a password

clean:
	docker compose down -v --remove-orphans
	docker system prune -af

# Start both postgres and keydb
deps:
	docker compose up -d postgres keydb

# Stop both postgres and keydb
deps-down:
	docker compose stop postgres keydb

# Start only postgres
postgres-up:
	docker compose up -d postgres

# Stop only postgres
postgres-down:
	docker compose stop postgres

# Start only keydb
keydb-up:
	docker compose up -d keydb

# Stop only keydb
keydb-down:
	docker compose stop keydb
