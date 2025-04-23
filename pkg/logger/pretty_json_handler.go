package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/fatih/color"
)

type prettyJSONHandler struct {
	out       io.Writer
	level     slog.Level
	attrGroup []string
	attrs     []slog.Attr
}

func newPrettyJSONHandler(w io.Writer, level slog.Level) *prettyJSONHandler {
	return &prettyJSONHandler{
		out:   w,
		level: level,
	}
}

func (h *prettyJSONHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *prettyJSONHandler) Handle(_ context.Context, r slog.Record) error {
	logMap := make(map[string]interface{})
	logMap["time"] = r.Time.Format(time.RFC3339)
	logMap["level"] = r.Level.String()
	logMap["msg"] = r.Message

	for _, attr := range h.attrs {
		logMap[attr.Key] = attr.Value.Any()
	}

	r.Attrs(func(attr slog.Attr) bool {
		logMap[attr.Key] = attr.Value.Any()
		return true
	})

	jsonData, err := json.MarshalIndent(logMap, "", "  ")
	if err != nil {
		return err
	}

	var colorFunc func(format string, a ...interface{}) string
	switch r.Level {
	case slog.LevelDebug:
		colorFunc = color.New(color.FgCyan).SprintfFunc()
	case slog.LevelInfo:
		colorFunc = color.New(color.FgGreen).SprintfFunc()
	case slog.LevelWarn:
		colorFunc = color.New(color.FgYellow).SprintfFunc()
	case slog.LevelError:
		colorFunc = color.New(color.FgRed).SprintfFunc()
	default:
		colorFunc = fmt.Sprintf
	}

	fmt.Fprintln(h.out, colorFunc("%s", jsonData))
	return nil
}

func (h *prettyJSONHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandler := *h
	newHandler.attrs = append(h.attrs, attrs...)
	return &newHandler
}

func (h *prettyJSONHandler) WithGroup(name string) slog.Handler {
	newHandler := *h
	newHandler.attrGroup = append(h.attrGroup, name)
	return &newHandler
}
