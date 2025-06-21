module github.com/orduro/pos-microservices/gateway

go 1.24.1

require (
	github.com/go-chi/chi/v5 v5.2.1
	github.com/go-chi/cors v1.2.1
	github.com/joho/godotenv v1.5.1
	github.com/orduro/pos-microservices/pkg v0.0.0-00010101000000-000000000000
)

replace github.com/orduro/pos-microservices/pkg => ../pkg
