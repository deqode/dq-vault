
export PGPASSWORD ?= $(POSTGRES_PASSWORD)
export DB_ENTRY ?= psql -h $(DB_HOST) -p 5432 -U $(POSTGRES_USER)

install-docker:
	@echo "Installing Docker"

	@sudo apt-get update

	@sudo apt-get install \
		apt-transport-https \
		ca-certificates \
		curl \
		software-properties-common -y

	@curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -

	@sudo add-apt-repository \
		"deb [arch=amd64] https://download.docker.com/linux/ubuntu \
		$$(lsb_release -cs) \
		stable"

	@sudo apt-get update

	@sudo apt-get --yes --no-install-recommends install docker-ce

	@sudo usermod --append --groups docker "$$USER"

	@sudo systemctl enable docker

	@echo "Waiting for Docker to start..."

	@sleep 3

	@sudo curl -L https://github.com/docker/compose/releases/download/1.22.0/docker-compose-`uname -s`-`uname -m` -o /usr/local/bin/docker-compose

	@sudo chmod +x /usr/local/bin/docker-compose
	@sleep 5
	@echo "Docker Installed successfully"

install-docker-if-not-already-installed:
	@if [ -z "$$(which docker)" ]; then\
		make install-docker;\
	fi

remove_stopped_containers:
	@docker-compose rm -v

down:
	@docker-compose down
	@docker-compose kill

build-all-docker-images:
	@docker-compose build --force-rm

set-up: install-docker-if-not-already-installed down remove_stopped_containers build-all-docker-images

dirty-up:
	@docker-compose up

run:	wait-for-postgres
	@vault server -config=/vault/config/config.hcl

wait-for-postgres:
	while ! nc -zv ${DB_HOST} 5432; do echo waiting for postgresql ..; sleep 5; done;

dev:
	@vault server -dev

.PHONY: run dev
