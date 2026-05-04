import { FormEvent } from "react";

type ChatInputProps = {
    value: string;
    placeholder: string;
    isLoading: boolean;
    onChange: (value: string) => void;
    onSubmit: (event: FormEvent<HTMLFormElement>) => void;
};

// Componente reutilizável do campo de mensagem.
//
// Ele não sabe nada sobre Agent, API, OpenWeather ou estado global.
// Ele só recebe valor, evento de mudança e evento de submit.
export function ChatInput({
    value,
    placeholder,
    isLoading,
    onChange,
    onSubmit,
}: ChatInputProps) {
    return (
        <form onSubmit={onSubmit} className="w-full">
            <div className="flex items-center gap-3 rounded-3xl border border-white/10 bg-[#151922] px-4 py-3 shadow-2xl shadow-black/30">
                <input
                    value={value}
                    onChange={(event) => onChange(event.target.value)}
                    placeholder={placeholder}
                    className="h-10 flex-1 bg-transparent text-sm text-zinc-100 outline-none placeholder:text-zinc-500"
                />

                <button
                    type="submit"
                    disabled={isLoading || !value.trim()}
                    className="flex h-10 w-10 items-center justify-center rounded-full bg-sky-500 text-lg font-medium text-white transition hover:bg-sky-400 disabled:cursor-not-allowed disabled:bg-zinc-700 disabled:text-zinc-400"
                >
                    →
                </button>
            </div>
        </form>
    );
}