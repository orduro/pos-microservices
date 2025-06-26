# this command generates swagger docs for the gateway
# make sure to run when updating docs or adding new docs 
# for new endpoints
swagger:
	@swag init \
	--generalInfo ./cmd/api/main.go \
	--output ./cmd/api/docs \
	--parseInternal

.PHONY: swagger
