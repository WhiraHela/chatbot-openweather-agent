import { tool } from "@langchain/core/tools";
import { z } from "zod";
import { fetchBackend } from "./api_client";

export const getCurrentWeatherTool = tool(
    async ({ city }: { city: string }) => {
        return fetchBackend(`/weather?city=${encodeURIComponent(city)}`);
    },
    {
        name: "get_current_weather",
        description:
            "Busca o clima atual de uma cidade. Use quando o usuário perguntar sobre temperatura atual, sensação térmica, umidade, vento ou condição climática agora.",
        schema: z.object({
            city: z.string().describe("Nome da cidade. Exemplo: São Paulo."),
        }),
    }
);

export const getForecastTool = tool(
    async ({ city, days }: { city: string; days?: number }) => {
        const safeDays = days && days > 0 ? days : 3;

        return fetchBackend(
            `/forecast?city=${encodeURIComponent(city)}&days=${safeDays}`
        );
    },
    {
        name: "get_weather_forecast",
        description:
            "Busca a previsão do tempo para os próximos dias de uma cidade. Use quando o usuário perguntar sobre amanhã, próximos dias, fim de semana, chuva futura ou previsão.",
        schema: z.object({
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

export const weatherTools = [getCurrentWeatherTool, getForecastTool];