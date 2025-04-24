package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"time"

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

	var parsed map[string]interface{}
	if err := json.Unmarshal(jsonData, &parsed); err != nil {
		return err
	}

	var buf bytes.Buffer
	buf.WriteString("{\n")

	first := true
	for k, v := range parsed {
		if !first {
			buf.WriteString(",\n")
		}
		first = false

		buf.WriteString("  ")
		buf.WriteString(h.keyColor.Sprintf(`"%s": `, k))

		switch k {
		case "level":
			if colorFunc, ok := h.levelColorMap[r.Level]; ok {
				buf.WriteString(colorFunc.Sprintf(`"%s"`, v))
			} else {
				buf.WriteString(h.valueColor.Sprintf(`"%s"`, v))
			}
		case "msg":
			underlined := color.New(color.Underline).Add(color.FgWhite)
			buf.WriteString(underlined.Sprintf(`"%s"`, v))
		default:
			valBytes, _ := json.Marshal(v)
			buf.WriteString(h.valueColor.Sprint(string(valBytes)))
		}
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
