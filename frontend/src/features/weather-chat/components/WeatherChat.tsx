"use client";

import { ChatInput } from "./ChatInput";
import { ChatMessage } from "./ChatMessage";
import { EmptyChat } from "./EmptyChat";
import { ThinkingIndicator } from "./ThinkingIndicator";
import { useWeatherChat } from "../hooks/useWeatherChat";

// Componente principal do chat.
//
// Ele funciona como "container visual":
// - usa o hook useWeatherChat para obter estado e ações
// - decide se mostra tela inicial ou conversa
// - renderiza os componentes menores
export function WeatherChat() {
    const {
        messages,
        input,
        setInput,
        isLoading,
        hasMessages,
        thinkingText,
        endRef,
        handleSubmit,
    } = useWeatherChat();

    return (
        <main className="min-h-screen bg-[#080a0d] text-white">
            <div className="mx-auto flex min-h-screen w-full max-w-4xl flex-col px-4">
                {!hasMessages ? (
                    <EmptyChat
                        input={input}
                        isLoading={isLoading}
                        onInputChange={setInput}
                        onSubmit={handleSubmit}
                    />
                ) : (
                    <>
                        <section className="flex-1 overflow-y-auto py-8">
                            <div className="mx-auto flex w-full max-w-3xl flex-col gap-5">
                                {messages.map((message, index) => (
                                    <ChatMessage key={index} message={message} />
                                ))}

                                {isLoading && (
                                    <ThinkingIndicator text={thinkingText} />
                                )}

                                <div ref={endRef} />
                            </div>
                        </section>

                        <footer className="sticky bottom-0 bg-[#080a0d]/95 py-5 backdrop-blur">
                            <div className="mx-auto w-full max-w-3xl">
                                <ChatInput
                                    value={input}
                                    isLoading={isLoading}
                                    onChange={setInput}
                                    onSubmit={handleSubmit}
                                    placeholder="Digite sua pergunta..."
                                />
                            </div>
                        </footer>
                    </>
                )}
            </div>
        </main>
    );
}