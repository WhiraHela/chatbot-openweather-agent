"use client";

import { FormEvent, useEffect, useMemo, useRef, useState } from "react";

type Message = {//formato de cada msg do chat
    role: "user" | "assistant";
    content: string;
};

const thinkingSteps = [
    "Analisando sua pergunta",
    "Escolhendo a melhor ferramenta",
    "Consultando os dados do clima",
    "Preparando a resposta",
];

export function WeatherChat() {
    const [messages, setMessages] = useState<Message[]>([]);//guarda todas as msg da conversa
    const [input, setInput] = useState("");//guarda texto q o usuario esta digitando no campo
    const [isLoading, setIsLoading] = useState(false);// sistema eesta esperando resposta do agente
    const [tick, setTick] = useState(0); //contador para alternar frases de fases de pensamento

    const endRef = useRef<HTMLDivElement | null>(null);// cria ref para div no final da conversa

    const hasMessages = messages.length > 0;

    useEffect(() => { //alterna frases de fase de pensamento no front
        //roda sempre que isLoading muda
        if (!isLoading) {
            setTick(0); //zera o efeito
            return;
        }

        const interval = setInterval(() => {
            setTick((prev) => prev + 1); //aumenta o tick a cada 450ms
        }, 450);//intervalo de 450ms

        return () => clearInterval(interval); //limpa intervalo quando o efeito for desmontado e isLoading mudar
    }, [isLoading]);


    
    useEffect(() => { //isLoading muda, pega a div final e "scrolla" ate ela
        endRef.current?.scrollIntoView({ behavior: "smooth" }); //(? -> optional chaining)
    }, [messages, isLoading]);



    const thinkingText = useMemo(() => { //monta p texto de pensamento
        const step = thinkingSteps[Math.floor(tick / 4) % thinkingSteps.length];
        const dots = ".".repeat((tick % 3) + 1);

        return `${step}${dots}`;
    }, [tick]);


    //roda quando o usuario evia o formulario, chamada em: <form onSubmit={handleSubmit}>
    async function handleSubmit(event: FormEvent<HTMLFormElement>) {
        event.preventDefault(); //cancela o evento padrao (reload da pagina)

        const userMessage = input.trim();

        if (!userMessage || isLoading) return; //se a msg estiver varia ou se ja esta carregando - retorna

        //atualiza o array de mensagens
        setMessages((prev) => [
            ...prev, //copia msgs antigas
            { //adiciona nova mensagem do usuario
                role: "user",
                content: userMessage,
            },
        ]);
        //^ estado em react deve ser tratado como imutavel, cria novo array para nao modificar array antigo (estado anterior)

        //limpa input e ativa loading
        setInput("");
        setIsLoading(true);

        //^ botao fica desabilitado, texto de pensamento aparece, tick comeca a aumentar

        try { //tenta enviar pergunta para rota interna do Next.js
            const response = await fetch("/api/chat", { //requisicao para POST /api/chat -> arquivo: frontend/src/app/api/chat/route.ts
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ //envia msg em json
                    message: userMessage,
                }),
            });

            const data = await response.json(); //le json retornado por /api/chat

            //adiciona nova msg do assistente no chat
            setMessages((prev) => [
                ...prev,
                {
                    role: "assistant",
                    content:
                        response.ok && data?.answer 
                        ? data.answer // se status.ok for true e existir data.answer usa data.answer existente
                        : data?.error || "Não foi possível processar sua mensagem.", // senao, use data.error - se data.error nao existir use mensagem padrao
                },
            ]);
        } 
        catch { //se fetch falhar -> erro de conexao
            setMessages((prev) => [
                ...prev,
                {
                    role: "assistant",
                    content:
                        "Não foi possível conectar ao assistente. Verifique se o backend está rodando.",
                },
            ]);
        } 
        finally { //finaliza loading -> tanto para sucesso quanto para erro (seta loading como false - destrava botao e input)
            setIsLoading(false);
        }
    }

    // ====JSX PRINCIPAL===

    return ( //retorna interface visual
        <main className="min-h-screen bg-[#080a0d] text-white">
        <div className="mx-auto flex min-h-screen w-full max-w-4xl flex-col px-4">
            {!hasMessages ? ( //se hasMessage for falso mostra tela inicial
            <section className="flex flex-1 flex-col items-center justify-center"> 
                <div className="mb-8 text-center">
                <h1 className="text-3xl font-semibold tracking-tight text-zinc-100 sm:text-4xl">
                    Como posso te ajudar hoje?
                </h1>

                <p className="mt-3 text-sm text-zinc-500">
                    Pergunte sobre clima atual, previsão ou chance de chuva.
                </p>
                </div>

                <form onSubmit={handleSubmit} className="w-full max-w-2xl">
                <div className="flex items-center gap-3 rounded-3xl border border-white/10 bg-[#151922] px-4 py-3 shadow-2xl shadow-black/30">
                    <input //input controlaado por react
                    value={input}
                    onChange={(event) => setInput(event.target.value)}
                    placeholder="Ex: Vai chover amanhã em São Paulo?"
                    className="h-10 flex-1 bg-transparent text-sm text-zinc-100 outline-none placeholder:text-zinc-500"
                    />

                    <button
                    type="submit" 
                    disabled={isLoading || !input.trim()} //botao fica desabilitado quando isLoading for true ou se input estiver vazio
                    className="flex h-10 w-10 items-center justify-center rounded-full bg-sky-500 text-lg font-medium text-white transition hover:bg-sky-400 disabled:cursor-not-allowed disabled:bg-zinc-700 disabled:text-zinc-400"
                    >
                    → 
                    </button>
                </div>
                </form>
            </section>
            ) 
            : ( //se hasMessages for verdadeiro mostra a conversa
            <>
                <section className="flex-1 overflow-y-auto py-8">
                <div className="mx-auto flex w-full max-w-3xl flex-col gap-5">
                    {messages.map((message, index) => {
                    const isUser = message.role === "user";

                    return (
                        <div
                        key={index}
                        className={`flex ${
                            isUser ? "justify-end" : "justify-start" // é usuario ? bolha vai para direita : bolha vai para esquerda 
                        }`}
                        >
                        <div //estilo das bolhas
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
                    })}

                    {isLoading && ( // Texto de pensamento - se isLoading for true, renderiza bloco ; se for false, nao renderiza nada
                    <div className="flex justify-start">
                        <div className="rounded-3xl px-5 py-3 text-sm text-zinc-500">
                        {thinkingText}
                        </div>
                    </div>
                    )} 

                    <div ref={endRef} /> 
                </div> 
                </section>

                <footer className="sticky bottom-0 bg-[#080a0d]/95 py-5 backdrop-blur">
                <form onSubmit={handleSubmit} className="mx-auto w-full max-w-3xl">
                    <div className="flex items-center gap-3 rounded-3xl border border-white/10 bg-[#151922] px-4 py-3 shadow-2xl shadow-black/30">
                    <input
                        value={input}
                        onChange={(event) => setInput(event.target.value)}
                        placeholder="Digite sua pergunta..."
                        className="h-10 flex-1 bg-transparent text-sm text-zinc-100 outline-none placeholder:text-zinc-500"
                    />

                    <button
                        type="submit"
                        disabled={isLoading || !input.trim()}
                        className="flex h-10 w-10 items-center justify-center rounded-full bg-sky-500 text-lg font-medium text-white transition hover:bg-sky-400 disabled:cursor-not-allowed disabled:bg-zinc-700 disabled:text-zinc-400"
                    >
                        →
                    </button>
                    </div>
                </form>
                </footer>
            </>
            )}
        </div>
        </main>
    );
}