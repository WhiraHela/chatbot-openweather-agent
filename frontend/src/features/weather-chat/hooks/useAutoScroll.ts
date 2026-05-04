import { RefObject, useEffect } from "react";

// Hook responsável por rolar a conversa automaticamente até o final.
//
// Recebe uma ref apontando para uma div invisível no fim da lista de mensagens.
// Sempre que as dependências mudam, ele chama scrollIntoView.
export function useAutoScroll(
    endRef: RefObject<HTMLDivElement | null>,
    dependencies: unknown[]
) {
    useEffect(() => {
        endRef.current?.scrollIntoView({
            behavior: "smooth",
        });
    }, dependencies);
}