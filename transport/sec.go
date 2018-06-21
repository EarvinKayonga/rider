package transport

import (
	"crypto/tls"
	"net/http"
)

// HTTP Header for security.
const (
	HSTS    = "max-age=63072000; preload"
	HSTSKey = "Strict-Transport-Security"

	XFrameDeny = "DENY"
	XFrameKey  = "X-Frame-Options"

	NoSniff  = "nosniff"
	XContent = "X-Content-Type-Options"

	JSONContentType = "application/json"
	XContentTypeKey = "Content-Type"
)

// secureHeaders adds headers to all the response coming for the server.
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(HSTSKey, HSTS)
		w.Header().Add(XFrameKey, XFrameDeny)
		w.Header().Add(XContent, NoSniff)

		w.Header().Add(XContentTypeKey, JSONContentType)

		next.ServeHTTP(w, r)
	})
}

// TLS configuration for https server.
// like accepted ciphers
var (
	TLSConfiguration = &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler))
)
