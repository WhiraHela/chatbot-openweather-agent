import { logger } from "../../../shared/logger/logger";

const API_URL = process.env.BACKEND_API_URL || "http://localhost:8080";

function mapBackendError(status: number, backendMessage?: string) {
    const normalizedMessage = backendMessage?.toLowerCase() || "";

    if (status === 404) {
        return {
            __app_error: true,
            status: 404,
            code: "CITY_NOT_FOUND",
            message: "Não encontrei essa cidade. Verifique o nome e tente novamente.",
            technicalMessage: backendMessage,
        };
    }

    if (status === 400) {
        return {
            __app_error: true,
            status: 400,
            code: "INVALID_WEATHER_REQUEST",
            message: "A solicitação enviada para o serviço de clima é inválida.",
            technicalMessage: backendMessage,
        };
    }

    if (
        status === 500 &&
        (normalizedMessage.includes("openweather_api_key") ||
            normalizedMessage.includes("openweather_current_url") ||
            normalizedMessage.includes("openweather_forecast_url") ||
            normalizedMessage.includes("não configurada") ||
            normalizedMessage.includes("nao configurada"))
    ) {
        return {
            __app_error: true,
            status: 500,
            code: "WEATHER_SERVICE_CONFIG_ERROR",
            message:
                "O serviço de clima não está configurado corretamente. Verifique as variáveis de ambiente do backend.",
            technicalMessage: backendMessage,
        };
    }

    if (status === 502) {
        return {
            __app_error: true,
            status: 502,
            code: "WEATHER_SERVICE_UNAVAILABLE",
            message:
                "Não foi possível consultar o serviço externo de clima no momento. Verifique a configuração da OpenWeather e tente novamente.",
            technicalMessage: backendMessage,
        };
    }

    return {
        __app_error: true,
        status: status || 500,
        code: "UNKNOWN_ERROR",
        message: "Erro interno ao consultar o serviço de clima.",
        technicalMessage: backendMessage,
    };
}

// Cliente HTTP responsável por chamar o backend Go.
// Recebe apenas o endpoint, por exemplo:
// /weather?city=São%20Paulo
// /forecast?city=São%20Paulo&days=3
export async function fetchBackend(endpoint: string) {
    const url = `${API_URL}${endpoint}`;

    await logger.info(`Iniciando chamada ao backend: ${url}`);

    try {
        const response = await fetch(url);

        const data = await response.json().catch(() => null);

        // IMPORTANTE:
        // fetch não lança erro automaticamente para HTTP 400/404/500/502.
        // Então precisamos verificar response.ok manualmente.
        // Aqui retornamos um JSON especial para o weatherAgent detectar depois.
        if (!response.ok) {
            const backendMessage =
                typeof data?.error === "string"
                    ? data.error
                    : "Erro desconhecido retornado pelo backend.";

            const mappedError = mapBackendError(response.status, backendMessage);

            await logger.error(
                `Backend retornou erro. Status: ${response.status}, code: ${mappedError.code}, endpoint: ${endpoint}, message: ${backendMessage}`
            );

            return JSON.stringify(mappedError);
        }

        await logger.info(
            `Backend respondeu com sucesso. Status: ${response.status}, endpoint: ${endpoint}`
        );

        return JSON.stringify(data);
    } catch (error) {
        await logger.critical(
            `Falha crítica ao chamar backend: ${
                error instanceof Error ? error.message : String(error)
            }`
        );

        return JSON.stringify({
            __app_error: true,
            status: 503,
            code: "BACKEND_UNAVAILABLE",
            message:
                "Não foi possível conectar ao backend de clima. Verifique se o servidor Go está rodando.",
            technicalMessage:
                error instanceof Error ? error.message : String(error),
        });
    }
}