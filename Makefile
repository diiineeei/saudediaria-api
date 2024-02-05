# Variaveis para especificar propriedades do sistema
PROJECT=adega-do-vilao
SERVICE=app-lojinha
VERSION=latest

auth:
	gcloud auth configure-docker -q

# Faz o build do container
build:
	docker build --platform linux/amd64 --tag gcr.io/$(PROJECT)/$(SERVICE):$(VERSION) .

# Faz  a publicação do container no container registry do google cloud
publish:
	docker push gcr.io/$(PROJECT)/$(SERVICE):$(VERSION)

# Faz o build do container e publica no container registry do google cloud
deploy: build publish

# Comando para parar e remover os container em execussão
stop:
	docker stop $(docker ps -a -q)
	docker rm $(docker ps -a -q)

# Comando para executar testes locais
test:
	go clean -cache
	go test ./...
