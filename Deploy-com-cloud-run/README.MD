# Objetivo: Desenvolver um sistema em Go que receba um CEP, identifica a cidade e retorna o clima atual (temperatura em graus celsius, fahrenheit e kelvin). Esse sistema deverá ser publicado no Google Cloud Run.

## Entrega:

* O código-fonte completo da implementação. 
* Testes automatizados demonstrando o funcionamento.
* Utilize docker/docker-compose para que possamos realizar os testes de sua aplicação.
* Deploy realizado no Google Cloud Run (free tier) e endereço ativo para ser acessado.


### Configuração da API do https://www.weatherapi.com/ 
Antes de rodar o projeto localmente, crie o arquivo .env, no diretório do projeto, substituindo SUA-CHAVE, pela sua chave da API.  
Conteúdo do .env:  
API_KEY=SUA-CHAVE

---

### O link do projeto no Cloud Run está no linnk: 
https://deploy-com-cloud-run-908929792116.us-central1.run.app/temperaturebycep

### Casos de uso:
#### Teste com um CEP existente
https://deploy-com-cloud-run-908929792116.us-central1.run.app/temperaturebycep?cep=01001000

#### Teste com CEP inexistente
https://deploy-com-cloud-run-908929792116.us-central1.run.app/temperaturebycep?cep=12345678

#### Teste com CEP inválido (CEP com menos de 8 dígitos)
https://deploy-com-cloud-run-908929792116.us-central1.run.app/temperaturebycep?cep=1234

#### Teste com CEP inválido (CEP com caracteres não numéricos)
https://deploy-com-cloud-run-908929792116.us-central1.run.app/temperaturebycep?cep=12ab5678


