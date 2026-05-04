import { logger } from "../../../shared/logger/logger";

const API_URL = process.env.BACKEND_API_URL || "http://localhost:8080";

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

        if (!response.ok) {
            await logger.error(
                `Backend retornou erro. Status: ${response.status}, endpoint: ${endpoint}`
            );

            return JSON.stringify({
                error: true,
                message:
                    data?.error ||
                    "Não foi possível buscar os dados no serviço de clima.",
            });
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

        throw error;
    }
}