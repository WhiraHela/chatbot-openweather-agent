import { useEffect, useMemo, useState } from "react";
import { thinkingSteps } from "../constants/thinkingSteps";

// Hook responsável apenas pela animação do texto de pensamento.
//
// Exemplo:
// "Analisando sua pergunta."
// "Analisando sua pergunta.."
// "Analisando sua pergunta..."
// "Escolhendo a melhor ferramenta..."
export function useThinkingText(isLoading: boolean) {
    const [tick, setTick] = useState(0);

    useEffect(() => {
        // Quando não está carregando, zera o contador da animação.
        if (!isLoading) {
            setTick(0);
            return;
        }

        // Enquanto está carregando, incrementa o contador a cada 450ms.
        const interval = setInterval(() => {
            setTick((prev) => prev + 1);
        }, 450);

        // Cleanup: remove o intervalo quando o componente desmontar
        // ou quando isLoading mudar.
        return () => clearInterval(interval);
    }, [isLoading]);

    const thinkingText = useMemo(() => {
        // Troca a frase a cada 4 ticks.
        const step = thinkingSteps[Math.floor(tick / 4) % thinkingSteps.length];

        // Alterna entre ".", ".." e "...".
        const dots = ".".repeat((tick % 3) + 1);

        return `${step}${dots}`;
    }, [tick]);

    return thinkingText;
}