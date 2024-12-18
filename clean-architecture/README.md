## Instruções para execução do projeto
### Iniciar os containers do Docker, do rabbitmq, e da aplicação com o comando:
```bash
docker-compose up --build -d
```

### Caso seja necessário refazer o build do docker-compose.yaml, rodar os comandos:
```bash
docker-compose down --volumes
sudo rm -r .docker/
docker-compose up --build -d
```
### Portas dos serviços:
    web server on port :8000
    gRPC server on port 50051
    GraphQL server on port 8080




## Testes:

### Testar o gRPC com evans:
```bash
evans -r repl
package pb
service OrderService
call ListOrders
```

### Testar GraphQL:
- Acessar playground: http://localhost:8080/
- rodar a query:
```
query queryOrders {
  listOrders{
    id
    Price
    Tax
    FinalPrice
  }
}
```

### Testar API REST:
- acessar api/list_orders.http
- ou diretamente o link: http://localhost:8000/list


---

### Minhas Anotações:
Comando para entrar no banco de dados:
```bash
docker exec -it mysql sh -c 'mysql -uroot -p orders'
```

Atualizar protofiles(gRPC):
```bash
protoc --go_out=. --go-grpc_out=. internal/infra/grpc/protofiles/order.proto 
```

Atualizar schema do graphQL:
```bash
go run github.com/99designs/gqlgen generate
```

### Criar a fila no RabbitMQ:
- Criar fila orders
- Entrar na fila orders e fazer o bind com amq.direct


### Checar logs dos containers:
```docker-compose logs -f migrations
docker-compose logs -f goapp
docker-compose logs -f mysql
docker-compose logs -f rabbitmq
```
