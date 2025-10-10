ORG    := yykhomenko
NAME   := mars
REPO   := ${ORG}/${NAME}
TAG    := ${REPO}:$(shell date '+%y%m%d')-$(shell git log -1 --pretty=format:'%h')
LATEST := ${REPO}:latest

include .env
export

update: ## Update dependencies
	go get -u ./...
	go mod tidy

lint: ## Run linters
	golangci-lint run --no-config --issues-exit-code=0 --timeout=10m \
    --disable-all --enable=deadcode  --enable=gocyclo --enable=revive --enable=varcheck \
    --enable=structcheck --enable=maligned --enable=errcheck --enable=dupl --enable=ineffassign \
    --enable=interfacer --enable=unconvert --enable=goconst --enable=gosec --enable=megacheck

run:
	go run ./...

test:	## Run tests
	go test -race -timeout 30s ./...

bench: ## Run benchmarks
	go test ./... -bench=. -benchmem

clean: ## Clean project
	find . -name '.DS_Store' -type f -delete

image: ## Build image
	docker build -t ${TAG} -t ${LATEST} .

push: ## Push image
	docker push ${TAG} && \
	docker push ${LATEST}

pull: ## Pull images
	docker compose pull

build: ## Build containers
	docker compose build

start: ## Create and start containers
	docker compose up -d

stop: ## Stop and remove containers
	docker compose down

restart: ## Restart containers
	make stop build start

status: ## Print containers status
	docker compose ps

log: ## Print log
	docker compose logs -f

send:
	curl -v \
	-H "Content-Type: application/json" \
	-d '{"from": "123", "to": "456", "text":"789"}' \
	-X POST http://localhost:8080/messages

test_mars:
	curl "http://localhost:8080/messages?from=777&to=380670000001&text=Hello"

bench_mars:
	~/src/go/wrkb/wrkb mars http://localhost:8080/messages\?from\=__RANDI64_700_777__\&to\=__RANDI64_380670000001_380670099999__\&text\=__RANDSTR_lettersdigits_16__

bench_hash:
	~/src/go/wrkb/wrkb hashes http://127.0.0.1:8082/hashes/__RANDI64_380670000001_380679999999__

test_sis:
	curl http://localhost:9001/subscribers/380670000001

bench_sis:
	~/src/go/wrkb/wrkb sis http://127.0.0.1:9001/subscribers/__RANDI64_380670000001_380679999999__

help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / \
  {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
