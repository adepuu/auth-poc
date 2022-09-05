pkgs          = $(shell go list ./... | grep -vE '(vendor|mock)')
NOW=$(shell date)
compose_file=./docker/docker-compose.yml
compose=docker-compose -f ${compose_file}
auth_service_binary=auth
user_service_binary=user
warehouse_service_binary=warehouse
current_dir=$(shell pwd)

docker-start:
	@echo "${NOW} STARTING CONTAINER..."
	@${compose} up -d --build

docker-run-auth:
	@echo "${NOW} BUILDING..."
	@cd ./svc/auth && go mod vendor && go build -gcflags="all=-N -l" -o ./../../bin/${auth_service_binary} ./main.go
	@echo "${NOW} RUNNING..."
	@docker exec -it auth /usr/local/bin/${auth_service_binary}

docker-run-user:
	@echo "${NOW} BUILDING..."
	@cd ./svc/user && go mod vendor && go build -gcflags="all=-N -l" -o ./../../bin/${user_service_binary} ./main.go
	@echo "${NOW} RUNNING..."
	@docker exec -it user /usr/local/bin/${user_service_binary}
