import { ChatMessage as ChatMessageType } from "../types/chat";

type ChatMessageProps = {
    message: ChatMessageType;
};

// Renderiza uma única bolha de mensagem.
//
// Se for usuário, alinha à direita.
// Se for assistente, alinha à esquerda.
export function ChatMessage({ message }: ChatMessageProps) {
    const isUser = message.role === "user";

    return (
        <div className={`flex ${isUser ? "justify-end" : "justify-start"}`}>
            <div
                className={`max-w-[82%] whitespace-pre-wrap rounded-3xl px-5 py-3 text-sm leading-relaxed shadow-lg ${
                    isUser
                        ? "bg-sky-200 text-slate-950 shadow-sky-950/10"
                        : "bg-sky-300/80 text-slate-950 shadow-sky-950/10"
                }`}
            >
                {message.content}
            </div>
        </div>
    );
}