package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type LogLevel string

// define constantes para niveis de log
const (
	INFO     LogLevel = "INFO"
	WARNING  LogLevel = "WARNING"
	ERROR    LogLevel = "ERROR"
	CRITICAL LogLevel = "CRITICAL"
)

// define constantes para cores de niveis de log
const (
	colorReset   = "\033[0m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorRed     = "\033[31m"
	colorPurple  = "\033[35m"
)

func getLogDir() string {
	logDir := os.Getenv("BACKEND_LOG_DIR")
	if logDir == "" {
		logDir = "logs"
	}

	return logDir
}

func getLogFilePath() string { // monta caminho de arquivo com nome backend-data.log
	today := time.Now().Format("2006-01-02")
	fileName := fmt.Sprintf("backend-%s.log", today)

	return filepath.Join(getLogDir(), fileName)
}

func getFormattedTime() string { //formata datetime para linhas de log
	return time.Now().Format("02/01/2006 15:04")
}


func colorizeLevel(level LogLevel) string {
	levelText := fmt.Sprintf("[%s]", level)

	switch level {
	case INFO:
		return colorGreen + levelText + colorReset
	case WARNING:
		return colorYellow + levelText + colorReset
	case ERROR:
		return colorRed + levelText + colorReset
	case CRITICAL:
		return colorPurple + levelText + colorReset
	default:
		return levelText
	}
}

func getCallerPath() string {
	_, file, _, ok := runtime.Caller(3)
	if !ok {
		return "unknown"
	}

	workingDir, err := os.Getwd()
	if err != nil {
		return file
	}

	relativePath, err := filepath.Rel(workingDir, file)
	if err != nil {
		return file
	}

	return strings.ReplaceAll(relativePath, "\\", "/")
}

func ensureLogFileStarted(filePath string) {
	info, err := os.Stat(filePath)

	if err == nil && info.Size() > 0 {
		return
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Println("Erro ao criar arquivo de log:", err)
		return
	}
	defer file.Close()

	startLine := fmt.Sprintf(
		"========== LOG INICIADO EM %s ==========\n",
		time.Now().Format("02/01/2006 15:04:05"),
	)

	if _, err := file.WriteString(startLine); err != nil {
		log.Println("Erro ao escrever início do log:", err)
	}
}

func write(level LogLevel, message string) {
	filePath := getLogFilePath()
	dirPath := filepath.Dir(filePath)

	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		log.Println("Erro ao criar diretório de logs:", err)
		return
	}

	ensureLogFileStarted(filePath)

	callerPath := getCallerPath()
	dateTime := getFormattedTime()

	fileLine := fmt.Sprintf(
		"[%s] [%s] %s - %s\n",
		callerPath,
		level,
		dateTime,
		message,
	)

	//retorna string formatada da linha de log
	consoleLine := fmt.Sprintf(
		"[%s] %s %s - %s",
		callerPath,
		colorizeLevel(level),
		dateTime,
		message,
	)

	//retorna linha formatada no stdout
	fmt.Println(consoleLine)


	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Println("Erro ao abrir arquivo de log:", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(fileLine); err != nil {
		log.Println("Erro ao escrever arquivo de log:", err)
	}
}


//cria funcoes para escrever log em niveis de logs diferentes (padronizado)

func Info(message string) {
	write(INFO, message)
}

func Warning(message string) {
	write(WARNING, message)
}

func Error(message string) {
	write(ERROR, message)
}

func Critical(message string) {
	write(CRITICAL, message)
}