import { WeatherChat } from "../components/chat/weather_chat";

//define a interface de uma rota (representa rota inicial /src/app)
export default function Home() {
    return <WeatherChat />; //pagina principal deve mostrar o componenete WeatherChat
}