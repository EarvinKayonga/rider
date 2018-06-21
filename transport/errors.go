package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/EarvinKayonga/rider/domain"
	"github.com/EarvinKayonga/rider/logging"
	"github.com/EarvinKayonga/rider/storage"
)

// Erroring centralize the error handling on the transport layer.
func Erroring(_ context.Context, w http.ResponseWriter, err error, logger logging.Logger) {
	switch err {
	case domain.ErrBikeInUse:
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "bike already in use",
		})
		if err != nil {
			logger.WithError(err).Info("while json encoding a error")
		}

	case storage.ErrBikeNotFound:
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "bike not found",
		})
		if err != nil {
			logger.WithError(err).Info("while json encoding a error")
		}

	default:
		w.WriteHeader(http.StatusInternalServerError)
		_, err := fmt.Fprint(w, "An expected error occured :(")
		if err != nil {
			logger.WithError(err).Info("while json encoding a error")
		}
	}
}
