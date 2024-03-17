include scripts/*.mk

dev-up:
	export $(cat deployments/development/.env | xargs) && docker-compose -f deployments/development/docker-compose.yaml up


dev-down:
	docker-compose -f deployments/development/docker-compose.yaml down

dev-clean:
	docker-compose -f deployments/development/docker-compose.yaml stop && docker-compose -f deployments/development/docker-compose.yaml rm -f


test:
	@go test -v ./...