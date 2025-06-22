package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
	FATAL
)

func (l LogLevel) String() string {
	return [...]string{"DEBUG", "INFO", "WARNING", "ERROR", "FATAL"}[l]
}

func LogLevelFromString(level string) LogLevel {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "WARNING":
		return WARNING
	case "ERROR":
		return ERROR
	case "FATAL":
		return FATAL
	default:
		return INFO // Уровень по умолчанию
	}
}

type LogEntry struct {
	Time             time.Time
	File             string
	Line             int
	Level            LogLevel
	Header           string
	Body             string
	Container        string
	AdditionalLabels []Label
}

type Label struct {
	Key   string
	Value string
}

var Logger *log.Logger
var logChannel chan LogEntry
var containerName string
var minLogLevel LogLevel
var fileInfoLogLevel LogLevel
var lokiURL string
var lokiUsername string
var lokiPassword string
var lokiCooldown time.Time
var duplicateLogs bool

func init() {

	lokiURL = os.Getenv("LOKI_URL")
	lokiUsername = os.Getenv("LOKI_USERNAME")
	lokiPassword = os.Getenv("LOKI_PASSWORD")
	duplicateLogs = os.Getenv("DUPLICATE_LOGS") == "ON"

	logDir := os.Getenv("LOG_DIR")
	file, err := os.OpenFile(logDir+"/"+getContainerName()+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	Logger = log.New(file, "", 0)
	logChannel = make(chan LogEntry, 100)

	logLevelEnv := os.Getenv("LOG_LEVEL")
	minLogLevel = LogLevelFromString(logLevelEnv)

	fileInfoLogLevelEnv := os.Getenv("FILE_INFO_DEBUG_LEVEL")
	fileInfoLogLevel = LogLevelFromString(fileInfoLogLevelEnv)

	go processLogs()
}

func getContainerName() string {
	if containerName == "" {
		containerName = os.Getenv("CONTAINER_NAME")
	}
	return containerName
}

func processLogs() {
	for entry := range logChannel {
		if time.Since(lokiCooldown) > 5*time.Second {
			err := SendToLoki(entry)
			if err != nil {
				lokiCooldown = time.Now()
				Logger.Printf(
					"[%s] [%s:%d] [%s] [%s] %s: %s\n",
					entry.Time.Format(time.RFC3339),
					"logger.go",
					0,
					"ERROR",
					entry.Container,
					"Failed to send log to Loki",
					err.Error(),
				)
			} else if !duplicateLogs {
				continue
			}

		}
		if entry.Level >= minLogLevel {
			Logger.Printf(
				"[%s] [%s:%d] [%s] [%s] %s: %s\n",
				entry.Time.Format(time.RFC3339),
				entry.File,
				entry.Line,
				entry.Level,
				entry.Container,
				entry.Header,
				entry.Body)
		}
	}
}

func Log(level LogLevel, header, body string, details ...Label) {
	entry := LogEntry{
		Time:             time.Now(),
		Level:            level,
		Header:           header,
		Body:             body,
		Container:        getContainerName(),
		AdditionalLabels: details,
	}

	if level >= fileInfoLogLevel {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			entry.File = file
			entry.Line = line
		} else {
			entry.File = "unknown"
			entry.Line = 0
		}
	}

	go func() {
		logChannel <- entry
	}()
}

func LogFatal(args ...int) { // Get log data from recovery, ues arg if need additional params
	if r := recover(); r != nil {
		depth := 2
		if len(args) > 0 {
			depth = args[0]
		}
		_, file, line, ok := runtime.Caller(depth)
		var entry LogEntry
		entry.Time = time.Now()
		entry.Level = FATAL
		entry.Header = "Panic occurred"
		if ok {
			entry.Body = fmt.Sprintf("Recovered from panic: %+v\nFile: %s, Line: %d", r, file, line)
		} else {
			entry.Body = fmt.Sprintf("Recovered from panic: %+v", r)
		}
		entry.Container = getContainerName()

		go func() {
			logChannel <- entry
		}()
		fmt.Fprintln(os.Stderr, entry.Body)
	}
}
