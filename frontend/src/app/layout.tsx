import "./globals.css"; //aplicar css global na aplicacao toda

//define metadados da pagina
export const metadata = {
    title: "Weather Agent",
    description: "Chat com agente de clima",
};

//retorna estrutur HTML base da aplicacao 
export default function RootLayout({ //componente react
                                        children,
                                        }: {
                                            children: React.ReactNode;
                                    })  //func recebe prop children (pode ser qualque conteudo react valido)
    {
        return (//retorna JSX q sera renderizado
            <html lang="pt-BR">
                <body>{children}</body>
            </html>
        );
    }