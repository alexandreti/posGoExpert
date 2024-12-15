Olá devs!
Agora é a hora de botar a mão na massa. Para este desafio, você precisará criar o usecase de listagem das orders.
Esta listagem precisa ser feita com:
- Endpoint REST (GET /order)
- Service ListOrders com GRPC
- Query ListOrders GraphQL
Não esqueça de criar as migrações necessárias e o arquivo api.http com a request para criar e listar as orders.

Para a criação do banco de dados, utilize o Docker (Dockerfile / docker-compose.yaml), com isso ao rodar o comando docker compose up tudo deverá subir, preparando o banco de dados.
Inclua um README.md com os passos a serem executados no desafio e a porta em que a aplicação deverá responder em cada serviço.

Listagem das orders:
- REST
- GRPC
- GraphQL

Checar: https://plataforma.fullcycle.com.br/courses/c2957fa4-1e88-4425-be86-5a17ad2664ca/346/197/177/conteudos?capitulo=177&conteudo=9693


1)Iniciar os containers do Docker e do rabbitmq com o comando:
docker-compose up -d
Obs, Caso os containers já existam, podem ser deletados antes da criação com o comando:
docker rm -f $(docker ps -a -q)

2)Iniciar o banco de dados com o comando:
make migrate

3)Criar a fila no RabbitMQ
Criar fila orders
Entrar na fila orders e fazer o bind com amq.direct

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