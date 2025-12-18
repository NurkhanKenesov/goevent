package logging

import (
    "encoding/json"
    "io"
    "os"
    "sync"
    "time"
)

type Logger struct {
    out io.Writer
    mu  sync.Mutex
}

var std = New(os.Stdout)

func New(w io.Writer) *Logger {
    return &Logger{out: w}
}

func Std() *Logger { return std }

func (l *Logger) log(level, message string, fields map[string]interface{}) {
    l.mu.Lock()
    defer l.mu.Unlock()
    entry := map[string]interface{}{
        "timestamp":  time.Now().UTC().Format(time.RFC3339Nano),
        "level":      level,
        "message":    message,
    }
    for k, v := range fields {
        entry[k] = v
    }
    b, _ := json.Marshal(entry)
    l.out.Write(append(b, '\n'))
}

func Info(msg string, fields map[string]interface{}) {
    std.log("info", msg, fields)
}

func Error(msg string, fields map[string]interface{}) {
    std.log("error", msg, fields)
}

func Warn(msg string, fields map[string]interface{}) {
    std.log("warn", msg, fields)
}
