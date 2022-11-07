build:
	cd app && CGO_ENABLED=0 go build -o native main.go && cd ..
	docker build -t songjiayang/native-histogram-demo:v0.1 .
up:
	docker-compose up -d
down:
	docker-compose down