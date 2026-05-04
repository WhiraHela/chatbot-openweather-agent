import { FormEvent, useRef, useState } from "react";
import { sendChatMessage } from "../services/chatApi";
import { ChatMessage } from "../types/chat";
import { useAutoScroll } from "./useAutoScroll";
import { useThinkingText } from "./useThinkingText";

// Hook principal do chat.
//
// Ele centraliza:
// - mensagens
// - input
// - loading
// - envio do formulário
// - chamada para /api/chat
// - scroll automático
// - texto de pensamento
export function useWeatherChat() {
    const [messages, setMessages] = useState<ChatMessage[]>([]);
    const [input, setInput] = useState("");
    const [isLoading, setIsLoading] = useState(false);

    // Ref para uma div invisível no final da conversa.
    // Usada para rolar automaticamente até a última mensagem.
    const endRef = useRef<HTMLDivElement | null>(null);

    const hasMessages = messages.length > 0;

    // Hook separado que calcula o texto animado enquanto o agente responde.
    const thinkingText = useThinkingText(isLoading);

    // Hook separado que mantém a conversa rolada até o final.
    useAutoScroll(endRef, [messages, isLoading]);

    async function handleSubmit(event: FormEvent<HTMLFormElement>) {
        event.preventDefault();

        const userMessage = input.trim();

        // Evita enviar mensagem vazia ou múltiplas mensagens enquanto carrega.
        if (!userMessage || isLoading) return;

        // Adiciona a mensagem do usuário na conversa.
        setMessages((prev) => [
            ...prev,
            {
                role: "user",
                content: userMessage,
            },
        ]);

        // Limpa o input e ativa estado de carregamento.
        setInput("");
        setIsLoading(true);

        try {
            const data = await sendChatMessage(userMessage);

            // Adiciona a resposta do assistente.
            setMessages((prev) => [
                ...prev,
                {
                    role: "assistant",
                    content:
                        data.answer ||
                        data.error ||
                        "Não foi possível processar sua mensagem.",
                },
            ]);
        } catch {
            // Erro de rede ou falha inesperada no fetch.
            setMessages((prev) => [
                ...prev,
                {
                    role: "assistant",
                    content:
                        "Não foi possível conectar ao assistente. Verifique se o backend está rodando.",
                },
            ]);
        } finally {
            // Garante que o loading seja desligado tanto em sucesso quanto em erro.
            setIsLoading(false);
        }
    }

    return {
        messages,
        input,
        setInput,
        isLoading,
        hasMessages,
        thinkingText,
        endRef,
        handleSubmit,
    };
}