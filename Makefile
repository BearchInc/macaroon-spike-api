VM_NAME=macaroons-vm

vm:
	@docker-machine create --driver virtualbox $(VM_NAME) || true
	eval "$(docker-machine env macaroons-vm)"
	docker-machine start $(VM_NAME)

build:
	docker build -t approvald .

run: rm-all vm build
	docker run --publish 6060:8080 --name approval-service --detach approvald

stop-all:
	docker stop --time=1 approval-service

rm-all: stop-all
	docker rm approval-service || true
	docker rmi approvald || true

destroy:
	docker-machine rm $(VM_NAME)