import { NextRequest, NextResponse } from "next/server";
import { askWeatherAgent } from "../../../lib/agent/weather-agent";
import { logger } from "../../../lib/logger/logger";

//ponte entre interface visual do chat e Agent(LangChain/OpenAI)
//cria endpoint interno   POST /api/chat

//funcao para responder requisicoes POST para /api/chat -> para Next.js conseguir usar como rota
export async function POST(req: NextRequest) { //disponibilizando funcao para Next.js exergar e usar
// uso de async pois outras funcionalidades (ler body da req, gravar log, chamar agent...) nao sao instantaneas - precisa esperar o Promisse dessas outras tarefas 
    try {
        const body = await req.json(); //leitura do json req - converte para obj javascript

        const message = body?.message; //se body existir, pegue body message

        
        if (!message || typeof message !== "string") {
            await logger.warning("Requisição /api/chat sem mensagem válida");

            return NextResponse.json(
                { error: "Mensagem obrigatória." },
                { status: 400 }
            );
        }

        await logger.info(`Mensagem recebida em /api/chat: ${message}`);

        //guarda resposta para msg enviada pelo usr
        const answer = await askWeatherAgent(message);
        await logger.info(`Resposta retornada por /api/chat: ${answer}`);

        //retorna resposta para msg enviada pelo usr
        return NextResponse.json({
            answer,
        });
    } catch (error) {
        await logger.error(
            `Erro em /api/chat: ${
                error instanceof Error ? error.message : String(error)
            }`
        );

        return NextResponse.json(
            { error: "Erro interno ao processar mensagem." },
            { status: 500 }
        );
    }
}