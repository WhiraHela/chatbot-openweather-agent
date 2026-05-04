import { logger } from "../logger/logger";

// 1. Recebe endpoint
// 2. Monta url completa
// 3. chama backend go com fetch
// 4. trata sucesso, erro HTTP e erro critico de conexao
// 5. retorna string json para a tool/agent

//pega url do dotenv, se n encontrar entra o fallback
const API_URL = process.env.BACKEND_API_URL || "http://localhost:8080";

//func assincrona q pode ser usada por outros arquivos
// async pois exerce acoes nao instantaneas
export async function fetchBackend(endpoint: string) {
    const url = `${API_URL}${endpoint}`; 

    await logger.info(`Iniciando chamada ao backend: ${url}`);
    try {
        const response = await fetch(url); //faz requisicao e espera o backend responder
        const data = await response.json().catch(() => null);//le json da reesposta e transforma em obs javascript

        if (!response.ok) { //trata erro de codigo http
            await logger.error(
                `Backend retornou erro. Status: ${response.status}, endpoint: ${endpoint}`
            );

            return JSON.stringify({ //retorna string json(stringfy) com erro padronizado
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
    } 
    catch (error) {
        await logger.critical(
            `Falha crítica ao chamar backend: ${
                error instanceof Error ? error.message : String(error)
            }`
        );

        throw error;
    }
}