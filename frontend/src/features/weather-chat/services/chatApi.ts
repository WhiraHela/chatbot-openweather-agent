import {
    ChatApiErrorResponse,
    ChatApiSuccessResponse,
    SendChatMessageResult,
} from "../types/chat";

export async function sendChatMessage(
    message: string
): Promise<SendChatMessageResult> {
    try {
        const response = await fetch("/api/chat", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                message,
            }),
        });

        const data = (await response.json()) as
            | ChatApiSuccessResponse
            | ChatApiErrorResponse;

        if (!response.ok) {
            const errorData = data as ChatApiErrorResponse;

            return {
                ok: false,
                status: response.status,
                code: errorData.code || "UNKNOWN_ERROR",
                message:
                    errorData.message ||
                    "Não foi possível processar sua mensagem.",
            };
        }

        const successData = data as ChatApiSuccessResponse;

        return {
            ok: true,
            answer: successData.answer,
        };
    } catch {
        return {
            ok: false,
            status: 503,
            code: "FRONTEND_NETWORK_ERROR",
            message:
                "Não foi possível conectar ao assistente. Verifique se o frontend está rodando corretamente.",
        };
    }
}