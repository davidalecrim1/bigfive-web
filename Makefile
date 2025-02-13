.PHONY: run local-run

run:
	@make generate
	@go run ./cmd/webui/main.go

local-run:
	@make generate
	air

generate:
	@templ generate