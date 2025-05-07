package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// LogLevel はログレベルを表す型
type LogLevel int

const (
	// DEBUG はデバッグレベル
	DEBUG LogLevel = iota
	// INFO は情報レベル
	INFO
	// WARN は警告レベル
	WARN
	// ERROR はエラーレベル
	ERROR
	// FATAL は致命的エラーレベル
	FATAL
)

// Logger はロガー構造体
type Logger struct {
	level  LogLevel
	logger *log.Logger
	file   *os.File
}

// NewLogger は新しいロガーを作成する
func NewLogger(level LogLevel, output io.Writer) *Logger {
	return &Logger{
		level:  level,
		logger: log.New(output, "", log.LstdFlags),
	}
}

// NewFileLogger はファイル出力のロガーを作成する
func NewFileLogger(level LogLevel, filePath string) (*Logger, error) {
	// ディレクトリが存在しない場合は作成
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	// ファイルを開く（追記モード）
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	return &Logger{
		level:  level,
		logger: log.New(file, "", log.LstdFlags),
		file:   file,
	}, nil
}

// Close はロガーを閉じる
func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// Debug はデバッグメッセージを記録する
func (l *Logger) Debug(format string, args ...interface{}) {
	if l.level <= DEBUG {
		l.log("DEBUG", format, args...)
	}
}

// Info は情報メッセージを記録する
func (l *Logger) Info(format string, args ...interface{}) {
	if l.level <= INFO {
		l.log("INFO", format, args...)
	}
}

// Warn は警告メッセージを記録する
func (l *Logger) Warn(format string, args ...interface{}) {
	if l.level <= WARN {
		l.log("WARN", format, args...)
	}
}

// Error はエラーメッセージを記録する
func (l *Logger) Error(format string, args ...interface{}) {
	if l.level <= ERROR {
		l.log("ERROR", format, args...)
	}
}

// Fatal は致命的エラーメッセージを記録し、プログラムを終了する
func (l *Logger) Fatal(format string, args ...interface{}) {
	if l.level <= FATAL {
		l.log("FATAL", format, args...)
		os.Exit(1)
	}
}

// log は実際にログを記録する内部メソッド
func (l *Logger) log(level string, format string, args ...interface{}) {
	// 呼び出し元の情報を取得
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	// ファイル名のみを取得
	file = filepath.Base(file)

	// タイムスタンプを取得
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// ログメッセージを整形
	msg := fmt.Sprintf(format, args...)
	logMsg := fmt.Sprintf("%s [%s] %s:%d - %s", timestamp, level, file, line, msg)

	// ログを記録
	l.logger.Println(logMsg)
}

// GetLogFilePath はログファイルのパスを取得する
func GetLogFilePath() (string, error) {
	// ホームディレクトリを取得
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// ログディレクトリのパス
	logDir := filepath.Join(homeDir, ".local", "share", "tui-chat", "logs")

	// 現在の日付を取得
	date := time.Now().Format("2006-01-02")

	// ログファイルのパス
	logPath := filepath.Join(logDir, fmt.Sprintf("chat-%s.log", date))

	return logPath, nil
}

// DefaultLogger はデフォルトのロガー
var DefaultLogger = NewLogger(INFO, os.Stdout)
