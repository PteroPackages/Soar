package logger

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Entry struct {
	color bool
	quiet bool
	time  string
	level string
	data  []string
	out   *os.File
}

func (e *Entry) Log() {
	border := e.level

	if e.color {
		switch e.level {
		case "INFO":
			border = "\x1b[34mINFO\x1b[0m"
		case "WARN":
			border = "\x1b[33mWARN\x1b[0m"
		case "ERROR":
			border = "\x1b[31mERROR\x1b[0m"
		case "":
			for _, line := range e.data {
				e.out.WriteString(fmt.Sprintf("%s\n", line))
			}

			return
		}
	}

	if e.level != "ERROR" && e.quiet {
		return
	}

	for _, line := range e.data {
		e.out.WriteString(fmt.Sprintf("%s: %s\n", border, line))
	}
}

func (e *Entry) Line(data string) *Entry {
	e.data = append(e.data, data)

	return e
}

func (e *Entry) WithError(err error) *Entry {
	return e.Line(err.Error())
}

func (e *Entry) WithTip(data string) *Entry {
	e.data = append(e.data, fmt.Sprintf("run '%s' for more information", data))

	return e
}

func (e *Entry) Format() []string {
	res := make([]string, len(e.data))

	for _, line := range e.data {
		res = append(res, fmt.Sprintf("[%s] %s: %s", e.time, e.level, line))
	}

	return res
}

type Logger struct {
	Color   bool
	Debug   bool
	Quiet   bool
	Persist bool

	entries []*Entry
}

func New() *Logger {
	return &Logger{
		Color:   true,
		Debug:   false,
		Quiet:   false,
		Persist: false,
		entries: []*Entry{},
	}
}

func (l *Logger) Line(data string) *Entry {
	e := Entry{
		color: l.Color,
		quiet: l.Quiet,
		time:  time.Now().Format(time.RFC822),
		level: "",
		data:  []string{data},
		out:   os.Stdout,
	}
	l.entries = append(l.entries, &e)

	return &e
}

func (l *Logger) Info(data string) *Entry {
	e := Entry{
		color: l.Color,
		quiet: l.Quiet,
		time:  time.Now().Format(time.RFC822),
		level: "INFO",
		data:  []string{data},
		out:   os.Stdout,
	}
	l.entries = append(l.entries, &e)

	return &e
}

func (l *Logger) Warn(data string) *Entry {
	e := Entry{
		color: l.Color,
		quiet: l.Quiet,
		time:  time.Now().Format(time.RFC822),
		level: "WARN",
		data:  []string{data},
		out:   os.Stdout,
	}
	l.entries = append(l.entries, &e)

	return &e
}

func (l *Logger) Error(data string) *Entry {
	e := Entry{
		color: l.Color,
		quiet: l.Quiet,
		time:  time.Now().Format(time.RFC822),
		level: "ERROR",
		data:  []string{data},
		out:   os.Stderr,
	}
	l.entries = append(l.entries, &e)

	return &e
}

func (l *Logger) Save() (string, error) {
	if len(l.entries) == 0 {
		return "", errors.New("no logs to save")
	}

	root, err := os.UserConfigDir()
	if err != nil {
		return "", errors.New("logs path is unavailable")
	}

	path := filepath.Join(root, "soar", "logs")
	info, err := os.Stat(path)
	if err != nil {
		return "", errors.New("logs path is unavailable")
	}

	if !info.IsDir() || info.Mode()&0o644 == 0 {
		return "", errors.New("logs path is invalid or unreachable")
	}

	name := strings.ReplaceAll(time.Now().Format(time.RFC3339), ":", "-")
	path = filepath.Join(path, name+".log")
	file, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	buf := strings.Builder{}
	for _, log := range l.entries {
		for _, line := range log.Format() {
			buf.WriteString(line + "\n")
		}
	}
	file.WriteString(buf.String())

	return path, nil
}
