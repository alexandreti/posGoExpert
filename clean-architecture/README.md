Olá devs!
Agora é a hora de botar a mão na massa. Para este desafio, você precisará criar o usecase de listagem das orders.
Esta listagem precisa ser feita com:
- Endpoint REST (GET /order)
- Service ListOrders com GRPC
- Query ListOrders GraphQL
Não esqueça de criar as migrações necessárias e o arquivo api.http com a request para criar e listar as orders.

Para a criação do banco de dados, utilize o Docker (Dockerfile / docker-compose.yaml), com isso ao rodar o comando docker compose up tudo deverá subir, preparando o banco de dados.
Inclua um README.md com os passos a serem executados no desafio e a porta em que a aplicação deverá responder em cada serviço.

1)Iniciar os containers do Docker e do rabbitmq com o comando:
docker-compose up -d

2)Criar o Banco de dados orders:
docker exec -it mysql sh -c 'mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS orders;"'
Digitar a senha do mysql: root

2)Iniciar o banco de dados com o comando:
make migrate

3) Iniciar o projeto com os comandos:
cd cmd/ordersystem/
go run main.go wire_gen.go

Portas dos serviços:
    web server on port :8000
    gRPC server on port 50051
    GraphQL server on port 8080


Comando para entrar no banco de dados:
    docker exec -it mysql sh -c 'mysql -uroot -p orders'
    Senha: root
    Comandos do banco de dados:
        show tables