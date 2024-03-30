# Configuração do Projeto

## Arquivo config.yaml

O arquivo config.yaml contém as configurações de limitação de taxa para endereços IP e tokens de API. Aqui está como você pode modificá-lo:

```yaml
rate_limit:
  ips:
    - ip: "192.168.1.1"
      limit: 10
    - ip: "192.168.1.2"
      limit: 2
      expiration: 5
      block: 10
  tokens:
    - token: "token123"
      limit: 2
      expiration: 5
      block: 10
```
- **limit**: número máximo de requisições permitidas antes de o sistema bloquear novas requisições do token/ip
- **expiration**: o tempo, em segundos, após o qual a contagem de requisições do token/ip é resetada
- **block**: o tempo, em segundos, pelo qual o token/ip é bloqueado após exceder o limite).

Quando uma dessas variaveis nao é definida, o sistema usa os valores padrão definidos no arquivo de configuração do docker-compose.yaml.

## Arquivo docker-compose.yaml

- No serviço rate_limiter, você pode configurar variáveis de ambiente como DEFAULT_IP_EXPIRATION_TIME, DEFAULT_IP_BLOCK_DURATION, etc., para definir limites padrão e comportamentos de limitação. 
  - **DEFAULT_IP_EXPIRATION_TIME**: o tempo, em segundos, após o qual a contagem de requisições do IP é resetada
  - **DEFAULT_TOKEN_EXPIRATION_TIME**: o tempo, em segundos, após o qual a contagem de requisições do token é resetada
  - **DEFAULT_IP_REQUEST_LIMIT**: número máximo de requisições permitidas antes de o sistema bloquear novas requisições do IP
  - **DEFAULT_TOKEN_REQUEST_LIMIT**: número máximo de requisições permitidas antes de o sistema bloquear novas requisições do token
  - **DEFAULT_IP_BLOCK_DURATION**: o tempo, em segundos, pelo qual o IP é bloqueado após exceder o limite
  - **DEFAULT_TOKEN_BLOCK_DURATION**: o tempo, em segundos, pelo qual o token é bloqueado após exceder o limite

## Executando o Projeto

```bash
docker-compose up --build
```

```bash
curl -X POST http://localhost:8080/ \
     -H "Content-Type: application/json" \
     -H "API_KEY: token123"
```