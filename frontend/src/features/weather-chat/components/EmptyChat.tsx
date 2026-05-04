import { FormEvent } from "react";
import { ChatInput } from "./ChatInput";

type EmptyChatProps = {
    input: string;
    isLoading: boolean;
    onInputChange: (value: string) => void;
    onSubmit: (event: FormEvent<HTMLFormElement>) => void;
};

// Tela inicial do chat, antes de existir qualquer mensagem.
export function EmptyChat({
    input,
    isLoading,
    onInputChange,
    onSubmit,
}: EmptyChatProps) {
    return (
        <section className="flex flex-1 flex-col items-center justify-center">
            <div className="mb-8 text-center">
                <h1 className="text-3xl font-semibold tracking-tight text-zinc-100 sm:text-4xl">
                    Como posso te ajudar hoje?
                </h1>

                <p className="mt-3 text-sm text-zinc-500">
                    Pergunte sobre clima atual, previsão ou chance de chuva.
                </p>
            </div>

            <div className="w-full max-w-2xl">
                <ChatInput
                    value={input}
                    isLoading={isLoading}
                    onChange={onInputChange}
                    onSubmit={onSubmit}
                    placeholder="Ex: Vai chover amanhã em São Paulo?"
                />
            </div>
        </section>
    );
}