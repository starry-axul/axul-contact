.PHONY: build

build:
	@echo "=> Building service"
	dir=contacts/get make build-file
	dir=contacts/create make build-file
	dir=contacts/update make build-file
	dir=contacts/getall make build-file
	dir=contacts/delete make build-file
	dir=contacts/alert make build-file
	sam build

local:
	make build
	sam local start-api --skip-pull-image --warm-containers EAGER --profile costamagna-terraform --env-vars env.json --docker-network appnet

build-file:
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/$(dir)/bootstrap		cmd/$(dir)/main.go

deploy:
	sam deploy --no-confirm-changeset --no-confirm-credentials --no-confirm-rollback --no-confirm-changeset-approval
