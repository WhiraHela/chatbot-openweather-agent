// Define os tipos usados pelos componentes e hooks do chat.

// Representa uma mensagem exibida na conversa.
export type ChatMessage = {
    // Define quem enviou a mensagem.
    // "user" = usuário
    // "assistant" = resposta do agente
    role: "user" | "assistant";

    // Texto da mensagem.
    content: string;
};

// Representa o formato esperado da resposta da rota /api/chat.
export type ChatApiResponse = {
    answer?: string;
    error?: string;
};