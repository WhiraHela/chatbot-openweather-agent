import { createAgent } from "langchain";
import { ChatOpenAI } from "@langchain/openai";
import { weatherTools } from "./tools";
import {
    AppError,
    appErrorFromSerialized,
    isAppError,
    isSerializedAppError,
} from "../weather-chat/types/appError";

const model = new ChatOpenAI({
    model: "gpt-4o-mini",
    temperature: 0.2,
});

export const weatherAgent = createAgent({
    model, //modelo do gpt
    tools: weatherTools, //tools disponiveis
    systemPrompt: //prommpt do agent 
        `
        Você é um assistente especializado exclusivamente em clima e previsão do tempo.

        Sua função é responder perguntas sobre clima usando obrigatoriamente as tools disponíveis.

        Você tem duas tools:

        1. get_current_weather
        Use quando o usuário perguntar sobre:
        - clima atual
        - temperatura agora
        - sensação térmica agora
        - umidade atual
        - vento agora
        - condição climática atual
        - se está chovendo agora

        2. get_weather_forecast
        Use quando o usuário perguntar sobre:
        - previsão do tempo
        - amanhã
        - próximos dias
        - fim de semana
        - se vai chover
        - clima futuro
        - previsão para 1 a 5 dias

        Regras obrigatórias:
        - Nunca invente dados climáticos.
        - Sempre use uma tool antes de responder perguntas sobre clima.
        - Se a pergunta for sobre agora, use get_current_weather.
        - Se a pergunta for sobre futuro, use get_weather_forecast.
        - Se o usuário não informar a cidade, peça a cidade antes de chamar qualquer tool.
        - Se a pergunta não for sobre clima ou previsão do tempo, responda educadamente que você só pode ajudar com clima e previsão do tempo.
        - Não use tool para perguntas fora do escopo de clima.
        - Responda sempre em português do Brasil.
        - Seja direto, claro e útil.
        - Use graus Celsius.
        - Quando houver chance de chuva, explique a probabilidade em porcentagem.
        - Não mencione detalhes técnicos da API, JSON, endpoint, backend ou OpenWeather para o usuário final.
        - Use os dados retornados pela tool como fonte de verdade.
        `,
});

// Tenta extrair um JSON de erro serializado de um conteúdo qualquer.
function tryExtractSerializedAppErrorFromText(text: string): AppError | null {
    try {
        const parsed = JSON.parse(text);

        if (isSerializedAppError(parsed)) {
            return appErrorFromSerialized(parsed);
        }
    } catch {
        // Se não for JSON puro, tentamos localizar o marcador no texto.
    }

    // Fallback:
    // Em alguns casos a mensagem pode vir com texto ao redor do JSON.
    const markerIndex = text.indexOf('"__app_error":true');

    if (markerIndex === -1) {
        return null;
    }

    const startIndex = text.lastIndexOf("{", markerIndex);
    const endIndex = text.indexOf("}", markerIndex);

    if (startIndex === -1 || endIndex === -1) {
        return null;
    }

    const possibleJson = text.slice(startIndex, endIndex + 1);

    try {
        const parsed = JSON.parse(possibleJson);

        if (isSerializedAppError(parsed)) {
            return appErrorFromSerialized(parsed);
        }
    } catch {
        return null;
    }

    return null;
}

// Percorre todas as mensagens retornadas pelo Agent.
// Se alguma tool retornou __app_error, interrompemos a resposta normal.
function throwIfAgentUsedAppError(messages: unknown[]) {
    for (const message of messages) {
        const content = (message as { content?: unknown })?.content;

        if (typeof content === "string") {
            const appError = tryExtractSerializedAppErrorFromText(content);

            if (appError) {
                throw appError;
            }
        }

        if (Array.isArray(content)) {
            for (const item of content) {
                if (typeof item === "string") {
                    const appError = tryExtractSerializedAppErrorFromText(item);

                    if (appError) {
                        throw appError;
                    }
                }

                if (
                    item &&
                    typeof item === "object" &&
                    "text" in item &&
                    typeof (item as { text?: unknown }).text === "string"
                ) {
                    const appError = tryExtractSerializedAppErrorFromText(
                        (item as { text: string }).text
                    );

                    if (appError) {
                        throw appError;
                    }
                }
            }
        }
    }
}

export async function askWeatherAgent(message: string) {
    try {
        const result = await weatherAgent.invoke({
            messages: [
                {
                    role: "user",
                    content: message,
                },
            ],
        });

        // PONTO PRINCIPAL:
        // Antes de pegar a última resposta do assistant, verifica se alguma tool
        // retornou erro técnico serializado.
        throwIfAgentUsedAppError(result.messages);

        const lastMessage = result.messages[result.messages.length - 1];

        if (typeof lastMessage.content === "string") {
            return lastMessage.content;
        }

        return JSON.stringify(lastMessage.content);
    } catch (error) {
        if (isAppError(error)) {
            throw error;
        }

        throw new AppError({
            code: "UNKNOWN_ERROR",
            status: 500,
            userMessage: "Erro interno ao processar sua mensagem.",
            technicalMessage:
                error instanceof Error ? error.message : String(error),
        });
    }
}