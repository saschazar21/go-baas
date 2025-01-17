package v1

import (
	"encoding/json"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/saschazar21/go-baas/booleans"
	"github.com/saschazar21/go-baas/db"
	"github.com/saschazar21/go-baas/errors"
)

func handleCreateBoolean(w http.ResponseWriter, r *http.Request) {
	b, err := booleans.ParseBoolean(r)

	if err != nil {
		httpErr, ok := err.(*errors.HTTPError)

		if !ok {
			httpErr = errors.NewHTTPError(http.StatusInternalServerError, &errors.INTERNAL_SERVER_ERROR)
		}

		httpErr.Write(w)
		return
	}

	var client *redis.Client
	if client, err = db.NewRedis(); err != nil {
		httpErr := errors.NewHTTPError(http.StatusInternalServerError, &errors.INTERNAL_SERVER_ERROR)
		httpErr.Write(w)
		return
	}

	defer client.Close()

	if err = b.Save(client, r.Context()); err != nil {
		httpErr, ok := err.(*errors.HTTPError)

		if !ok {
			httpErr = errors.NewHTTPError(http.StatusInternalServerError, &errors.INTERNAL_SERVER_ERROR)
		}

		httpErr.Write(w)
		return
	}

	res := booleans.CreateBooleanResponse(b)

	w.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(w).Encode(res); err != nil {
		httpErr := errors.NewHTTPError(http.StatusInternalServerError, &errors.INTERNAL_SERVER_ERROR)
		httpErr.Write(w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func HandleBooleans(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)

		httpErr := errors.NewHTTPError(http.StatusMethodNotAllowed, &errors.METHOD_NOT_ALLOWED_ERROR)
		httpErr.Write(w)
		return
	}

	params := r.URL.Query()
	params.Del("id")

	r.URL.RawQuery = params.Encode()

	handleCreateBoolean(w, r)
}
