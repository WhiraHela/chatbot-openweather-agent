# OpenWeather Chatbot (Go + Next.js + LangChain)

Aplicação fullstack para consultar clima e previsão do tempo por meio de um chat. O backend foi feito em Go com Gin e expõe endpoints REST para clima atual e forecast. O frontend foi feito em Next.js com TypeScript e usa um Agent com LangChain/OpenAI para interpretar a pergunta do usuário e escolher qual ferramenta chamar.

Exemplos de perguntas:

```txt
Qual o clima atual em São Paulo?
Vai chover amanhã em Ribeirão Preto?
Como estará o tempo nos próximos 3 dias em Curitiba?
```

Usei LangChain porque ele facilita a criação de Agents com tools. Nesse projeto, o Agent decide se deve buscar clima atual ou previsão futura, chama a ferramenta correta e usa o backend como fonte dos dados. Também escolhi LangChain por já ter experiência prévia com a ferramenta, o que ajudou na organização do fluxo entre LLM, tools e backend.



## Tecnologias

### Backend

* Go
* Gin
* OpenWeather API
* godotenv
* Logger próprio

### Frontend

* Next.js
* TypeScript
* React
* Tailwind CSS
* LangChain JS
* OpenAI API
* Zod



## Como rodar localmente

### Pré-requisitos

* Go
* Node.js + npm
* Chave da OpenWeather
* Chave da OpenAI



## Configuração do backend

Entre na pasta do backend:

```bash
cd backend
```

Copie o arquivo de exemplo:

```bash
cp .env.example .env
```

Edite `backend/.env`:

```env
OPENWEATHER_API_KEY=SUA_CHAVE_DA_OPENWEATHER
OPENWEATHER_CURRENT_URL=https://api.openweathermap.org/data/2.5/weather
OPENWEATHER_FORECAST_URL=https://api.openweathermap.org/data/2.5/forecast
SERVER_ADDRESS=localhost:8080
```



## Configuração do frontend

Entre na pasta do frontend:

```bash
cd frontend
```

Copie o arquivo de exemplo:

```bash
cp .env.local.example .env.local
```

Edite `frontend/.env.local`:

```env
BACKEND_API_URL=http://localhost:8080
OPENAI_API_KEY=SUA_CHAVE_DA_OPENAI
AGENT_LOG_DIR=logs
```



## Rodando o backend

```bash
cd backend
go mod download
go run .
```

Servidor esperado:

```txt
http://localhost:8080
```



## Rodando o frontend

Em outro terminal:

```bash
cd frontend
npm install
npm run dev
```

Acesse:

```txt
http://localhost:3000
```



## Endpoints do backend

```http
GET /health
GET /weather?city=São%20Paulo
GET /forecast?city=São%20Paulo&days=3
```

Exemplos:

```bash
curl "http://localhost:8080/health"
```

```bash
curl "http://localhost:8080/weather?city=Sao%20Paulo"
```

```bash
curl "http://localhost:8080/forecast?city=Sao%20Paulo&days=3"
```



## Rota interna do frontend

O frontend usa a rota:

```http
POST /api/chat
```

Exemplo:

```bash
curl -i -X POST "http://localhost:3000/api/chat" \
  -H "Content-Type: application/json" \
  -d '{"message":"qual o clima atual em sao paulo"}'
```

Resposta de sucesso:

```json
{
  "answer": "O clima atual em São Paulo é..."
}
```

Resposta de erro:

```json
{
  "error": true,
  "code": "WEATHER_SERVICE_UNAVAILABLE",
  "message": "Não foi possível consultar o serviço externo de clima no momento."
}
```



## Fluxo básico

```txt
Usuário
↓
Frontend Next.js
↓
POST /api/chat
↓
Agent LangChain/OpenAI
↓
Tool escolhida
↓
Backend Go/Gin
↓
OpenWeather API
↓
Backend retorna JSON limpo
↓
Agent gera resposta final
↓
Chat exibe a resposta
```



## Logs

```txt
backend/logs
frontend/logs
```

Os logs registram chamadas, erros controlados, falhas externas e respostas processadas.



## Tratamento de erros

O projeto separa respostas normais do Agent de erros reais da aplicação.

Respostas normais retornam `200`, por exemplo:

* pergunta fora do escopo
* pergunta sem cidade
* resposta normal de clima

Erros reais retornam status HTTP adequado, por exemplo:

* `404` para cidade não encontrada
* `500` para configuração ausente
* `502` para erro na OpenWeather
* `503` para backend indisponível



## Plano de testes

No navegador, abra:

```txt
F12 → Network → Fetch/XHR → chat
```

Confira sempre:

```txt
Headers → Status Code
Payload → mensagem enviada
Response → JSON retornado
```



## Testes funcionais e de erro

| Teste                        | Entrada/configuração               | Status esperado | Resposta esperada                    |
| ---------------------------- | ------------------------------------ | --------------: | ------------------------------------ |
| Health check                 | `GET /health`                      |             200 | API rodando                          |
| Clima atual válido          | `qual o clima atual em sao paulo`  |             200 | `answer`                           |
| Forecast válido             | `vai chover amanhã em sao paulo?` |             200 | `answer`                           |
| Fora do escopo               | `quem é o presidente do brasil?`  |             200 | Agent informa que só responde clima |
| Sem cidade                   | `qual o clima atual?`              |             200 | Agent pede a cidade                  |
| Cidade inexistente           | `qual o clima em xyzabczzzz?`      |             404 | `CITY_NOT_FOUND`                   |
| Backend desligado            | parar backend e perguntar clima      |             503 | `BACKEND_UNAVAILABLE`              |
| OpenWeather key ausente      | `OPENWEATHER_API_KEY=`             |             500 | `WEATHER_SERVICE_CONFIG_ERROR`     |
| OpenWeather key inválida    | chave inválida no backend           |             502 | `WEATHER_SERVICE_UNAVAILABLE`      |
| URL da OpenWeather incorreta | path incorreto na URL                |             502 | `WEATHER_SERVICE_UNAVAILABLE`      |
| OpenAI key ausente           | `OPENAI_API_KEY=`                  |             500 | `OPENAI_API_KEY_MISSING`           |



## Testes com curl

### Backend

```bash
curl -i "http://localhost:8080/health"
```

```bash
curl -i "http://localhost:8080/weather?city=Sao%20Paulo"
```

```bash
curl -i "http://localhost:8080/forecast?city=Sao%20Paulo&days=3"
```

```bash
curl -i "http://localhost:8080/weather?city=xyzabczzzz"
```



## Testes de responsividade

Use o DevTools do navegador:

```txt
F12 → Toggle device toolbar → Responsive
```

Tamanhos sugeridos:

| Tipo           |     Tamanho |
| -------------- | ----------: |
| Mobile pequeno |   360 x 640 |
| Mobile médio  |   390 x 844 |
| Tablet         |  768 x 1024 |
| Notebook       |  1366 x 768 |
| Desktop        | 1920 x 1080 |

Checklist:

```txt
[ ] Não existe scroll horizontal.
[ ] O input fica visível.
[ ] O botão de enviar não fica cortado.
[ ] Mensagens longas quebram linha corretamente.
[ ] A tela inicial fica centralizada.
[ ] Em desktop, o chat não fica largo demais.
```

Teste rápido no console:

```js
document.documentElement.scrollWidth > window.innerWidth
```

Esperado:

```txt
false
```



## Decisões técnicas

O backend foi separado do frontend para deixar a aplicação mais próxima de um cenário fullstack real. O Go/Gin ficou responsável por consultar a OpenWeather e devolver dados limpos. O Next.js ficou responsável pela interface, pela rota `/api/chat` e pela integração com o Agent.

O LangChain foi escolhido para organizar o Agent e as tools. Ele adiciona complexidade, mas facilita o roteamento entre clima atual e forecast, além de deixar mais claro quando a LLM deve usar uma ferramenta.

O tratamento de erro foi separado da resposta conversacional. Isso evita que erro técnico, como OpenWeather fora do ar ou chave ausente, vire uma resposta normal do chat com status `200`.



## Trade-offs

* Usar LangChain facilita o Agent, mas adiciona uma camada a mais de abstração.
* Separar backend e frontend melhora a organização, mas exige rodar dois serviços localmente.
* Ter logger próprio ajuda no entendimento do fluxo, mas em produção eu usaria uma solução mais robusta.
* Tratar erro fora da resposta da LLM deixa o sistema mais confiável, mas exige mais código de controle.



## O que eu faria diferente em produção

Em produção, eu melhoraria principalmente:

* gerenciamento de segredos com serviço próprio do provedor
* correlation ID por requisição
* rate limit
* cache para reduzir chamadas à OpenWeather
* testes unitários e de integração
* monitoramento com Sentry, Grafana, Datadog ou similar
* pipeline de CI/CD
* validação mais rígida das respostas externas



## Status do projeto

Projeto desenvolvido como desafio fullstack + AI, com foco em:

* integração entre Go e Next.js
* Agent com tools
* consumo de API externa
* tratamento de erros
* logs
* organização de código
* documentação técnica

