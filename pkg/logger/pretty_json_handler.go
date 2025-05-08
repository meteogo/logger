package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"

	"github.com/fatih/color"
)

type prettyJSONHandler struct {
	out           io.Writer
	level         slog.Level
	attrGroup     []string
	attrs         []slog.Attr
	keyColor      *color.Color
	valueColor    *color.Color
	levelColorMap map[slog.Level]*color.Color
}

func newPrettyJSONHandler(w io.Writer, level slog.Level) *prettyJSONHandler {
	return &prettyJSONHandler{
		out:        w,
		level:      level,
		keyColor:   color.New(color.FgHiWhite),
		valueColor: color.New(color.FgHiBlack),
		levelColorMap: map[slog.Level]*color.Color{
			slog.LevelDebug: color.New(color.FgMagenta),
			slog.LevelInfo:  color.New(color.FgBlue),
			slog.LevelWarn:  color.New(color.FgYellow),
			slog.LevelError: color.New(color.FgRed),
		},
	}
}

func (h *prettyJSONHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *prettyJSONHandler) Handle(_ context.Context, r slog.Record) error {
	logMap := make(map[string]interface{})
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

	var parsed map[string]interface{}
	if err := json.Unmarshal(jsonData, &parsed); err != nil {
		return err
	}

	var buf bytes.Buffer
	dateStr := r.Time.Format("02.01.2006T15:04:05")
	levelStr := r.Level.String()
	msgStr := r.Message
	levelColor, ok := h.levelColorMap[r.Level]
	if !ok {
		levelColor = h.valueColor
	}
	underlined := color.New(color.Underline).Add(color.FgWhite)
	prefix := fmt.Sprintf("[%s]~[%s]: %s\n", underlined.Sprintf("%s", dateStr), levelColor.Sprintf("%s", levelStr), underlined.Sprintf("%s", msgStr))
	buf.WriteString(prefix)

	if len(parsed) == 0 {
		fmt.Fprint(h.out, buf.String())
		return nil
	}

	buf.WriteString("{\n")
	first := true
	for k, v := range parsed {
		if !first {
			buf.WriteString(",\n")
		}
		first = false

		buf.WriteString("  ")
		buf.WriteString(h.keyColor.Sprintf(`"%s": `, k))

		valBytes, _ := json.Marshal(v)
		buf.WriteString(h.valueColor.Sprint(string(valBytes)))
	}
	buf.WriteString("\n}")

	fmt.Fprintln(h.out, buf.String())
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
