package middleware

type ContextKey string

const (
	ContextKeyRequestID = ContextKey("request_id")
	RequestIDUnknown    = "unknown"

	LogKeyStatus   = "status"
	LogKeyLatency  = "latency"
	LogKeyID       = "id"
	LogKeyMethod   = "method"
	LogKeyURI      = "uri"
	LogKeyHost     = "host"
	LogKeyRemoteIP = "remote_ip"
)
