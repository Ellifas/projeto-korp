.RECIPEPREFIX := >

PROJECT_NAME=projeto-korp
COMPOSE=docker compose
ANSIBLE=ansible-playbook -i ansible/inventory.ini ansible/playbook.yml

.PHONY: help test run build up down restart logs ps metrics ansible clean

help:
> @echo "Comandos disponíveis:"
> @echo "  make test      - Executa testes Go"
> @echo "  make run       - Executa aplicação localmente"
> @echo "  make build     - Faz build da imagem Docker"
> @echo "  make up        - Cria rede e sobe containers"
> @echo "  make down      - Remove containers"
> @echo "  make restart   - Reinicia ambiente"
> @echo "  make logs      - Exibe logs dos containers"
> @echo "  make ps        - Lista containers"
> @echo "  make metrics   - Consulta métrica korp_service_up"
> @echo "  make ansible   - Provisiona tudo via Ansible"
> @echo "  make clean     - Remove containers, volumes e imagem local"

test:
> go test ./...

run:
> go run ./cmd/server

build:
> docker build -t http-server-projeto-korp:latest .

up:
> docker network create korp-network 2>/dev/null || true
> $(COMPOSE) up -d --build

down:
> $(COMPOSE) down

restart: down up

logs:
> $(COMPOSE) logs -f

ps:
> $(COMPOSE) ps

metrics:
> curl "http://localhost:9090/api/v1/query?query=korp_service_up"

ansible:
> ansible-galaxy collection install -r ansible/requirements.yml
> $(ANSIBLE) -K

clean:
> $(COMPOSE) down -v --remove-orphans
> docker image rm http-server-projeto-korp:latest 2>/dev/null || true
