swagger:
	swag init -g main.go -d ./cmd/api,./internal\transport\http\handlers --parseDependency --parseInternal -o ./internal/transport/http/docs
