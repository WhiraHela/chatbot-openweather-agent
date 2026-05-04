import { appendFile, mkdir, stat } from "node:fs/promises";
import path from "node:path";

//define niveis de log
type LogLevel = "INFO" | "WARNING" | "ERROR" | "CRITICAL";

//define cores para cada nivel de log
const COLORS: Record<LogLevel, string> = {
    INFO: "\x1b[32m",
    WARNING: "\x1b[33m",
    ERROR: "\x1b[31m",
    CRITICAL: "\x1b[35m",
};

//cor da parte da mensagem da linha de log
const RESET = "\x1b[0m";


function getLogDir() { //retorna dir da pasta de logs
    return process.env.AGENT_LOG_DIR || "logs";
}

function getLogFilePath() { //retona caminho do arquivo de logs
    const today = new Date().toISOString().split("T")[0];

    return path.join(process.cwd(), getLogDir(), `agent-${today}.log`);
}

//formata datetime para cada linha de log
function formatDateTime() {
    const now = new Date();

    const day = String(now.getDate()).padStart(2, "0");
    const month = String(now.getMonth() + 1).padStart(2, "0");
    const year = now.getFullYear();

    const hours = String(now.getHours()).padStart(2, "0");
    const minutes = String(now.getMinutes()).padStart(2, "0");

    return `${day}/${month}/${year} ${hours}:${minutes}`;
}

//formatacao de cores para cada parte da linha de log
function colorizeLevel(level: LogLevel) {
    return `${COLORS[level]}[${level}]${RESET}`;
}

//retorna o arquivo que chamou o logger
function getCallerPath() {
    const stack = new Error().stack;    //retorna pilha de chamadas 
                                        // stack mostra quais funcoes foram chamadas (arquivo linha e coluna) - usa como trace das funcoes chamadas

    if (!stack) {
        return "unknown"; // caso n encontre o arquivo retorna unknow
    }

    const stackLines = stack.split("\n"); //quebra string em array de linhas
    const callerLine = stackLines[4] || stackLines[3] || ""; //tenta pegar a linha do arquivo que chamou o logger

    const match = // 2 opcoes de regex para pegar caminho da funcao que chama o logger
        callerLine.match(/\((.*):\d+:\d+\)/) ||
        callerLine.match(/at (.*):\d+:\d+/);

    if (!match) {
        return "unknown";
    }

    const absolutePath = match[1];

    return path //retorna caminho relativo
        .relative(process.cwd(), absolutePath)
        .replaceAll("\\", "/");
}


//garante que arquivo de logger tenha linha inicial
async function ensureLogStarted(filePath: string) {
    try { //verifica se o arquivo existe e se tem linha inicial
        const fileStat = await stat(filePath);

        if (fileStat.size > 0) {
            return;
        }
    } catch {
    // Se o arquivo ainda não existe, ele será criado abaixo.
    }
    //se nao tiver arquivo, cria abaixo e adiciona linha inicial no arquivo de log
    const startLine = `========== LOG INICIADO EM ${formatDateTime()} ==========\n`;

    await appendFile(filePath, startLine, "utf-8");
}

//funcao de escrita do logger
async function write(level: LogLevel, message: string) {
    // pega caminhos (caminho de arquivo e diretorio)
    const filePath = getLogFilePath();
    const dirPath = path.dirname(filePath);

    //tenta encontrar aqrquivo, se nao existir, cria diretorio, arquivo e adiciona linha inicial 
    await mkdir(dirPath, { recursive: true });
    await ensureLogStarted(filePath);

    //pega caminho da funcao que chamou o logger e pega datetime
    const callerPath = getCallerPath();
    const dateTime = formatDateTime();

    //define formato de linha do log
    const fileLine = `[${callerPath}] [${level}] ${dateTime} - ${message}\n`;

    const consoleLine = `[${callerPath}] ${colorizeLevel(
        level
    )} ${dateTime} - ${message}`;

    console.log(consoleLine); //imprime no terminal

    await appendFile(filePath, fileLine, "utf-8"); //adiciona linha de log no arquivo logger
}


//funcoes prontas para escrever log para cada nivel (padronizado)
export const logger = {
    info(message: string) {
        return write("INFO", message);
    },

    warning(message: string) {
        return write("WARNING", message);
    },

    error(message: string) {
        return write("ERROR", message);
    },

    critical(message: string) {
        return write("CRITICAL", message);
    },
};