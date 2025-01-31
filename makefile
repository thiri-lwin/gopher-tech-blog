dev:
	go run ./cmd/server/main.go
docker:
	docker build -t gopher-tech-blog .
	docker run -p 8080:8080 -e PORT=8080 -e DB_USER=<user> -e DB_PASSWORD=<password> gopher-tech-blog