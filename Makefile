build:
	cd app && CGO_ENABLED=0 go build -o client_golang main.go && cd ..
	docker build -t client_golang:v0.1 .
up:
	docker-compose up -d
down:
	docker-compose down