package logger

type contextKey string

const (
	contextKeyRequestID contextKey = contextKey("requestID")
	contextKeyEnvType   contextKey = contextKey("envType")
)
