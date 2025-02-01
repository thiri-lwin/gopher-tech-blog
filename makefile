dev:
	go run ./cmd/server/main.go
docker:
	docker build -t gopher-tech-blog .
	docker run -p 8080:8080 -e PORT=8080 -e DB_URI=<postgresql_uri> -e REDIS_ADDR=<addr> gopher-tech-blog
