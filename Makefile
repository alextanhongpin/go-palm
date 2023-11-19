include .env
export


run:
	@go run main.go


doc:
	@go run cmd/documenter/documenter.go < $(f)


tools:
	@go run cmd/tool/main.go
#palm -p="<EOF
#your text
#EOF"
