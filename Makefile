run:
	go run ./cmd/web

swag:
	swag init -g ./cmd/web/main.go

docker-build:
	docker compose build

docker-run:
	docker compose up --build

docker-run-dev:
	docker compose up --build --detach

docker-stop:
	docker compose down

docker-restart:
	docker compose down && docker compose up --build --detach

build-exec:
	go build -o ./bin/web ./cmd/web

run-exec:
	./bin/web
