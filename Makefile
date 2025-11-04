up:
	docker compose up --build 
down: 
	docker compose down -v
tests:
	go test ./internal/middleware -v 
	go test ./internal/handler -v