include .env
export


run:
	@go run main.go


doc:
	@go run cmd/documenter/documenter.go < $(f)
