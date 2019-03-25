test: down
	docker-compose -f docker/test.compose.yml up --build test
	make down

gen: down
	docker-compose -f docker/generate.compose.yml up --build generate
	make down

run: down
	docker-compose up -d --build api
	make logs-api



## Development helpers

dev: restart-api logs-api

restart-api:
	docker-compose kill api && docker-compose rm -v -f api && docker-compose up -d --no-deps --build api

build-api:
	docker-compose build api

ps:
	docker-compose ps
up: down
	docker-compose up -d
down:
	docker-compose -f docker/generate.compose.yml down && \
	docker-compose -f docker/test.compose.yml down && \
	docker-compose down

logs:
	docker-compose logs -f --tail=100
logs-api:
	docker-compose logs -f --tail=100 api

