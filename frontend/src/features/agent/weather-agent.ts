import { createAgent } from "langchain";
import { ChatOpenAI } from "@langchain/openai";
import { weatherTools } from "./tools";



const model = new ChatOpenAI({
    model: "gpt-4o-mini",
    temperature: 0.2, //resposta mais controlada e menos criativa, mais precisao para agent de clima
});


//cria o agent
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

//expor para poder sr usado por outros arquivos
export async function askWeatherAgent(message: string) {
    const result = await weatherAgent.invoke({ //chama agent com invoke e manda mensagem (role e content)
        messages: [
            {
                role: "user",
                content: message,
            },
        ],
    });

    const lastMessage = result.messages[result.messages.length - 1]; // array com mensagens geradas durante execucao

    if (typeof lastMessage.content === "string") { //valida se conteudo e string, se for retrona string
        return lastMessage.content; 
    }

    return JSON.stringify(lastMessage.content); // se conteuno nao for string retorna JSON
}