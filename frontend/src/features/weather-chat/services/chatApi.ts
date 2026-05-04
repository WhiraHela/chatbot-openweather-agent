import { ChatApiResponse } from "../types/chat";

// Chama a rota interna do Next.js responsável por conversar com o Agent.
// Essa função é usada pelo hook useWeatherChat, ou seja, pela interface do chat.
export async function sendChatMessage(message: string): Promise<ChatApiResponse> {
    const response = await fetch("/api/chat", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            message,
        }),
    });

    const data = (await response.json()) as ChatApiResponse;

    if (!response.ok) {
        return {
            error: data?.error || "Não foi possível processar sua mensagem.",
        };
    }

    return data;
}