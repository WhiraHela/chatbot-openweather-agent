import { getCurrentWeatherTool, getForecastTool } from "./src/lib/agent/tools";

async function main() {
    console.log("Testando tool de clima atual...");

    const currentWeather = await getCurrentWeatherTool.invoke({
        city: "São Paulo",
    });

    console.log(currentWeather);

    console.log("\nTestando tool de previsão...");

    const forecast = await getForecastTool.invoke({
        city: "São Paulo",
        days: 3,
    });

    console.log(forecast);
}

main().catch((error) => {
    console.error("Erro ao testar tools:", error);
});