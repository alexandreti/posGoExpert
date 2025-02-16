# Entrega:
- O código-fonte completo da implementação.
- Documentação explicando como rodar o projeto em ambiente dev.
- Utilize docker/docker-compose para que possamos realizar os testes de sua aplicação.

---
# Instruções para rodar o projeto:

### Baixar o código fonte:
https://github.com/alexandreti/posGoExpert/tree/main/observabilidade

### Configuração:
- Criar o arquivo .env no raiz do projeto, com a API do http://api.weatherapi.com/v1/current.json
- Conteúdo do arquivo .env:<br>
API_KEY=KEY

onde KEY corresponde a sua chave da API.

### Rodar o docker-compose
- Suba todos os serviços do docker-compose.yaml

### Execute as chamadas da API.
- Execute as chamadas da API para o serviço A: api/servicoa.http

### Acesse o zipkin
Para conferir os traces, acesse o link do zipkin: http://localhost:9412/
