.PHONY: run local-run

generate:
	@templ generate

run:
	@make generate
	@docker compose up --build -d