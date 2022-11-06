build:
	docker build -t client_golang:v0.1 .
up:
	docker-compose up -d
down:
	docker-compose down