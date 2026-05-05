import ReactMarkdown from "react-markdown";
import remarkBreaks from "remark-breaks";
import remarkGfm from "remark-gfm";
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
                {isUser ? (
                    message.content
                ) : (
                    <ReactMarkdown
                        remarkPlugins={[remarkGfm, remarkBreaks]}
                        components={{
                            p: ({ children }) => (
                                <p className="mb-3 last:mb-0">{children}</p>
                            ),
                            ul: ({ children }) => (
                                <ul className="mb-3 list-disc pl-6 last:mb-0">
                                    {children}
                                </ul>
                            ),
                            ol: ({ children }) => (
                                <ol className="mb-3 list-decimal pl-6 last:mb-0">
                                    {children}
                                </ol>
                            ),
                            li: ({ children }) => (
                                <li className="mb-1 last:mb-0">{children}</li>
                            ),
                            code: ({ children }) => (
                                <code className="rounded bg-slate-900/10 px-1 py-0.5 text-[0.85em]">
                                    {children}
                                </code>
                            ),
                            blockquote: ({ children }) => (
                                <blockquote className="border-l-4 border-slate-700/40 pl-3 italic">
                                    {children}
                                </blockquote>
                            ),
                            a: ({ children, href }) => (
                                <a
                                    className="underline underline-offset-4"
                                    href={href}
                                    target="_blank"
                                    rel="noreferrer"
                                >
                                    {children}
                                </a>
                            ),
                            strong: ({ children }) => (
                                <strong className="font-semibold">
                                    {children}
                                </strong>
                            ),
                        }}
                    >
                        {message.content}
                    </ReactMarkdown>
                )}
            </div>
        </div>
    );
}
