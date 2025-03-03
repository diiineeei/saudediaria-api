# APILojinha - Backend para lojas online

APILojinha é um sistema backend robusto desenvolvido em Go (Golang), projetado para facilitar a criação, gerenciamento e controle de lojas online. Com capacidade para suportar lojas de grande porte.

## Características

- **Escala**: Projetado para suportar grande escala.
- **Desempenho**: Otimizado para alta performance, garantindo tempos de resposta rápidos mesmo sob carga pesada.
- **Segurança**: Implementa as melhores práticas de segurança, garantindo a proteção dos dados dos participantes e a integridade do sorteio.
- **Flexibilidade**: API RESTful para fácil integração com frontends e outros sistemas.
- **Transparência**: Mecanismos integrados para garantir a transparência e justiça dos sorteios.

## Tecnologias Utilizadas

- **Linguagem**: Go (Golang)
- **Banco de Dados**: PostgreSQL/MySQL (Escolha conforme a necessidade)
- **Cache**: Redis (Para gerenciamento de estado em alta velocidade)
- **Autenticação**: JWT (JSON Web Tokens) para uma API segura
- **Logging**: Logrus/Zap para logging estruturado e eficiente

## Como Executar

- Instalação:

Clone o repositório do projeto:
> git clone https://github.com/diiineeei/ApiLojinha.git

Acesse o diretório do projeto:
> cd ApiLojinha

Instale as dependências do projeto:
> go mod download

Crie o arquivo de configuração do banco de dados:
> cp .env.example .env

Edite o arquivo .env e configure os detalhes do banco de dados.

- Inicie o servidor:

> go run main.go