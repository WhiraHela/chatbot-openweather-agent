type ThinkingIndicatorProps = {
    text: string;
};

// Renderiza o texto cinza exibido enquanto o agente está processando.
//
// Exemplo:
// "Consultando os dados do clima..."
export function ThinkingIndicator({ text }: ThinkingIndicatorProps) {
    return (
        <div className="flex justify-start">
            <div className="rounded-3xl px-5 py-3 text-sm text-zinc-500">
                {text}
            </div>
        </div>
    );
}