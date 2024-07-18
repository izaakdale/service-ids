run:
	USE_REDIS=true REDIS_ENDPOINT=redis://localhost:6379 REDIS_SERVER_PORT=8080 go run main.go