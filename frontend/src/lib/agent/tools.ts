import { tool } from "@langchain/core/tools";
import { z } from "zod"; //definir e validar formatos de dados
import { fetchBackend } from "./api_client";
import { logger } from "../logger/logger";


// O QUE CADA TOOL FAZ
    // 1. codifica cidade
    // 2. monta endpoint
    // 3. chama backend
    // 4. retorna resultado



// tool de clima atual
export const getCurrentWeatherTool = tool(
    //recebe um objeto city e extrai a propriedade desse objeto
    // uso de async pois chama backend (chamada http), qui usa o fetchBackend
    async ({ city }: { city: string }) => {
        await logger.info(`Tool get_current_weather selecionada para cidade: ${city}`);
        const endpoint = `/weather?city=${encodeURIComponent(city)}`;

        await logger.info(`Endpoint montado para clima atual: ${endpoint}`);
        const result = await fetchBackend(endpoint);

        await logger.info(`Tool get_current_weather finalizada com sucesso`);
        return result;
    },
    { // manual de instrucoes para o agent
        name: "get_current_weather",
        description:
            "Busca o clima atual de uma cidade. Use quando o usuário perguntar sobre temperatura atual, sensação térmica, umidade, vento ou condição climática agora.",
        schema: z.object({
            city: z.string().describe("Nome da cidade. Exemplo: São Paulo."),
        }),
    }
);

// tool de previsao (/forecast) para n dias e para cidade especificada
export const getForecastTool = tool(
    async ({ city, days }: { city: string; days?: number }) => {
        // se days existir e for mais que 0 use days, caso contrario use 3
        const safeDays = days && days > 0 ? days : 3;
        await logger.info(
            `Tool get_weather_forecast selecionada para cidade: ${city}, dias: ${safeDays}`
        );
        
        const endpoint = `/forecast?city=${encodeURIComponent(
            city
        )}&days=${safeDays}`;

        await logger.info(`Endpoint montado para previsão: ${endpoint}`);
        const result = await fetchBackend(endpoint);

        await logger.info(`Tool get_weather_forecast finalizada com sucesso`);

        return result;
    },
    {// descreve a tool para o agent
        name: "get_weather_forecast",
        description:
            "Busca a previsão do tempo para os próximos dias de uma cidade. Use quando o usuário perguntar sobre amanhã, próximos dias, fim de semana, chuva futura ou previsão.",
        schema: z.object({//entrada da tool precisa ser um obj com city e days
            city: z.string().describe("Nome da cidade. Exemplo: São Paulo."),
            days: z
                .number()
                .min(1)
                .max(5)
                .default(3)
                .describe("Número de dias da previsão, entre 1 e 5."),
        }),
    }
);

//exporta lista de tools para o agente escolher qual usar
export const weatherTools = [getCurrentWeatherTool, getForecastTool];