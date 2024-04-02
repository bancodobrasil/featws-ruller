include .env
export

lint:
	go install golang.org/x/lint/golint@latest
	@make run-on-our-code-directories ARGS="golint"

build:generate-swagger
	go build -o featws-ruller

run:build
	./featws-ruller

test:
	@make run-on-our-code-directories ARGS="go test -v"

run-on-our-code-directories:
	@echo "${ARGS} <our-code-directories>"
	@make our-code-directories | xargs -n 1 $(ARGS)
our-code-directories:
	@go list ./... | grep -v /docs
verify:test
	make lint

generate-swagger:
#   Install swag on https://github.com/swaggo/swag
	swag i

deps:
	go mod tidy

	