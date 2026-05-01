const API_URL = process.env.BACKEND_API_URL || "http://localhost:8080";

export async function fetchBackend(endpoint: string) {
    const response = await fetch(`${API_URL}${endpoint}`);

    const data = await response.json().catch(() => null);

    if (!response.ok) {
        return JSON.stringify({
            error: true,
            message:
                data?.error ||
                "Não foi possível buscar os dados no serviço de clima.",
        });
    }

    return JSON.stringify(data);
}