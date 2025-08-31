package logger

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"os"
)

type PrettyWriter struct {
	writer io.Writer
}

func (pw *PrettyWriter) Write(p []byte) (n int, err error) {
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, p, "", "  ")
	if err != nil {
		return 0, err
	}
	return pw.writer.Write(append(prettyJSON.Bytes(), '\n'))
}

func New() {
	//pw := &PrettyWriter{writer: os.Stdout}
	pw := os.Stdout
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}
	handler := slog.NewJSONHandler(pw, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
