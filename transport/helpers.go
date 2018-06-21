package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/handlers"

	"github.com/EarvinKayonga/rider/configuration"
)

// GetPaginationArguments extract limit and cursor from request.
func GetPaginationArguments(req *http.Request) (string, int64) {
	if req.URL == nil {
		return "", 20
	}

	queries := req.URL.Query()
	cursorID := queries.Get("cursor")

	limit, err := strconv.ParseInt(queries.Get("limit"), 10, 64)
	if err != nil {
		return cursorID, 20
	}

	return cursorID, limit
}

func health(_ context.Context, metadata Metadata) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		_ = json.NewEncoder(w).Encode(metadata)
	}
}

// NewServer sets up HTTP the server.
func NewServer(ctx context.Context, conf configuration.Server, router http.Handler) (*http.Server, error) {
	return &http.Server{
		Addr:    conf.String(),
		Handler: handlers.RecoveryHandler()(router),

		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}, nil
}
