# OpenWeather Chatbot (Go + Next.js)

Projeto fullstack com backend em Go (Gin) que expõe endpoints de clima/forecast e frontend em Next.js com agente LangChain/OpenAI para conversar sobre clima.

## Requisitos

- Go (versão definida em `backend/go.mod`)
- Node.js + npm (compatíveis com a versão do Next.js usada em `frontend/package.json`)
- Chave da OpenWeather
- Chave da OpenAI (para o agente)

## Variáveis de ambiente

### Backend (`backend/.env`)

Crie ou edite `backend/.env` com:

```
OPENWEATHER_API_KEY=COLOQUE_SUA_CHAVE_AQUI
OPENWEATHER_CURRENT_URL=https://api.openweathermap.org/data/2.5/weather
OPENWEATHER_FORECAST_URL=https://api.openweathermap.org/data/2.5/forecast
SERVER_ADDRESS=localhost:8080
```

### Frontend (`frontend/.env.local`)

Crie ou edite `frontend/.env.local` com:

```
BACKEND_API_URL=http://localhost:8080
OPENAI_API_KEY=COLOQUE_SUA_CHAVE_AQUI
AGENT_LOG_DIR=logs
```

Observações:

- `BACKEND_API_URL` deve apontar para o endereço do backend.
- `AGENT_LOG_DIR` é relativo à pasta do frontend.

## Como rodar localmente

### 1. Subir o backend

```
cd backend
go mod download

go run .
```

Saída esperada: servidor em `http://localhost:8080` (ou o valor de `SERVER_ADDRESS`).

### 2. Subir o frontend

Em outro terminal:

```
cd frontend
npm install
npm run dev
```

Acesse `http://localhost:3000`.

## Endpoints do backend

- `GET /health`
- `GET /weather?city=São%20Paulo`
- `GET /forecast?city=São%20Paulo&days=3`

Exemplos via curl:

```
curl "http://localhost:8080/health"

curl "http://localhost:8080/weather?city=Sao%20Paulo"

curl "http://localhost:8080/forecast?city=Sao%20Paulo&days=3"
```

## Logs

- Backend grava logs em `backend/logs`.
- Frontend grava logs em `frontend/logs` (configurável via `AGENT_LOG_DIR`).

## Problemas comuns

- Erro `OPENWEATHER_API_KEY não configurada`: verifique `backend/.env`.
- Erro `OPENWEATHER_CURRENT_URL` ou `OPENWEATHER_FORECAST_URL`: confira as URLs no `.env` do backend.
- Erro da OpenAI: verifique `OPENAI_API_KEY` em `frontend/.env.local`.
- Backend fora do ar: confirme se `go run .` está em execução e o `SERVER_ADDRESS` está correto.

## Estrutura geral

- `backend/`: API em Go (Gin)
- `frontend/`: Next.js + LangChain + OpenAI
- `systemDesign_openweather_chatbot.drawio`: diagrama do sistema
