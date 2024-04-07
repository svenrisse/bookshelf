include .env
export $(shell sed 's/=.*//' .env)

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run/api: run the cmd/api application
.PHONY: run/api 
run/api:
	air -- -db-dsn=${DB_DSN} -smtp-username=${SMTP_USERNAME} -smtp-password=${SMTP_PASSWORD}

## db: connect to db using psql
.PHONY: db 
db:
	psql ${DB_DSN}

## db/migrations/new name =$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations/ -database ${DB_DSN} up

## db/migrations/down: apply all down migrations
.PHONY: db/migrations/down
db/migrations/down: confirm
	@echo 'Running down migrations...'
	migrate -path ./migrations/ -database ${DB_DSN} down

## db/migrations/version: get migration version
.PHONY: db/migrations/version
db/migrations/version: 
	@echo 'Checking migration version'
	migrate -path ./migrations/ -database ${DB_DSN} version

## db/migrations/force: force to migration version
.PHONY: db/migrations/force
db/migrations/force: 
	@echo 'forcing migration version ${version}'
	migrate -path ./migrations/ -database ${DB_DSN} force ${version}



# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: tidy and vendor dependencies and format, vet and test all code
.PHONY: audit
audit: vendor
	@echo 'Formatting code...'
	go fmt ./...
	@echo "Formatting SWAG comments..."
	swag fmt 
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor

# ==================================================================================== #
# BUILD
# ==================================================================================== #
#

## build/api: build the cmd/api application
.PHONY: build/api
build/api: build/swagger
	@echo 'Building cmd/api...'
	go build -ldflags='-s' -o=./bin/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/api ./cmd/api

## build/swagger: generate the swagger docs
.PHONY: build/swagger
build/swagger: 
	swag init --parseDependency -g cmd/api/main.go

# ==================================================================================== #
# PRODUCTION
# ==================================================================================== #

## production/setup: connect to prod as root and run setup.sh
.PHONY: production/setup
production/setup:
	rsync -rP --delete ./remote/setup/01.sh root@${PROD_IP}:/root
	ssh -t root@${PROD_IP} "bash 01.sh"

## production/connect: connect to the production server
.PHONY: production/connect
production/connect:
	ssh ${PROD_USER}@${PROD_IP}

## production/deploy/api: deploy the api to production
.PHONY: production/deploy/api
production/deploy/api:
	rsync -P ./bin/linux_amd64/api ${PROD_USER}@${PROD_IP}:~/bookshelf/
	rsync -rP --delete ./migrations ${PROD_USER}@${PROD_IP}:~/bookshelf/
	rsync -P ./remote/production/Caddyfile ${PROD_USER}@${PROD_IP}:~/bookshelf/
	rsync -P ./remote/production/api.service ${PROD_USER}@${PROD_IP}:~/bookshelf/
	ssh -t ${PROD_USER}@${PROD_IP} '\
		migrate -path ~/bookshelf/migrations -database $$DB_DSN up \
		&& sudo mv ~/bookshelf/api.service /etc/systemd/system/ \
		&& sudo systemctl enable api \
		&& sudo systemctl restart api \
		&& sudo mv ~/bookshelf/Caddyfile /etc/caddy/ \
		&& sudo systemctl reload caddy \
	'
