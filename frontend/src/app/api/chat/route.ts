import { NextRequest, NextResponse } from "next/server";
import { askWeatherAgent } from "../../../features/agent/weatherAgent";
import { logger } from "../../../shared/logger/logger";
import { AppError, isAppError } from "../../../features/weather-chat/types/appError";

export async function POST(req: NextRequest) {
    try {
        const openaiApiKey = process.env.OPENAI_API_KEY;

        if (!openaiApiKey) {
            await logger.error("OPENAI_API_KEY não configurada no ambiente");

            return NextResponse.json(
                {
                    error: true,
                    code: "OPENAI_API_KEY_MISSING",
                    message: "Configure uma chave de API válida da OpenAI.",
                },
                { status: 500 }
            );
        }

        const body = await req.json();
        const message = body?.message;

        if (!message || typeof message !== "string") {
            await logger.warning("Requisição /api/chat sem mensagem válida");

            return NextResponse.json(
                {
                    error: true,
                    code: "MESSAGE_REQUIRED",
                    message: "Mensagem obrigatória.",
                },
                { status: 400 }
            );
        }

        await logger.info(`Mensagem recebida em /api/chat: ${message}`);

        const answer = await askWeatherAgent(message);

        await logger.info(`Resposta retornada por /api/chat: ${answer}`);

        return NextResponse.json(
            {
                answer,
            },
            { status: 200 }
        );
    } catch (error) {
        if (isAppError(error)) {
            await logger.error(
                `Erro controlado em /api/chat. Status: ${error.status}, code: ${error.code}, message: ${
                    error.technicalMessage || error.userMessage
                }`
            );

            return NextResponse.json(
                {
                    error: true,
                    code: error.code,
                    message: error.userMessage,
                },
                { status: error.status }
            );
        }

        await logger.error(
            `Erro em /api/chat: ${
                error instanceof Error ? error.message : String(error)
            }`
        );

        const unknownError = new AppError({
            code: "UNKNOWN_ERROR",
            status: 500,
            userMessage: "Erro interno ao processar mensagem.",
            technicalMessage:
                error instanceof Error ? error.message : String(error),
        });

        return NextResponse.json(
            {
                error: true,
                code: unknownError.code,
                message: unknownError.userMessage,
            },
            { status: unknownError.status }
        );
    }
}