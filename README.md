# CodePix | Imersão Full Stack && Full Cycle

> School of Net/Code.Education - Projeto Full Cycle Developer
> Repositório original • [codeedu/imersao-fullstack-fullcycle](https://github.com/codeedu/imersao-fullstack-fullcycle)

## Microsserviço CodePix

Tem o objetivo de intermediar todas as transações bancárias realizada nesse ecossistema fictício.

### Execução

#### Pré-requisitos

- Docker (e nada mais)

#### Comandos

1º passo, no terminal, pela raiz do diretório /codepix:
`docker-compose up -d`

2º passo, no terminal, acesse o container da aplicação:
`docker exec -t codepix_app bash`

3º passo, no terminal da aplicação codepix:
`go run cmd/codepix/main.go`

#### Recursos utilizados no docker-compose do CodePix

- PostgreSQL
- PGAdmin
- Apache Kafka
- Kafka Control Center
- Confluent Control Center
- ZooKeeper
- E a aplicação principal
