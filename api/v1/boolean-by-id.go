package v1

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/saschazar21/go-baas/booleans"
	"github.com/saschazar21/go-baas/db"
	"github.com/saschazar21/go-baas/errors"
)

func handleDeleteBooleanById(w http.ResponseWriter, r *http.Request, id string) {
	client, err := db.NewRedis()
	if err != nil {
		httpErr := errors.NewHTTPError(http.StatusInternalServerError, &errors.INTERNAL_SERVER_ERROR)
		httpErr.Write(w)
		return
	}

	if err := booleans.DeleteBoolean(client, r.Context(), id); err != nil {
		httpErr, ok := err.(*errors.HTTPError)

		if !ok {
			httpErr = errors.NewHTTPError(http.StatusInternalServerError, &errors.INTERNAL_SERVER_ERROR)
		}

		httpErr.Write(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleGetBooleanById(w http.ResponseWriter, r *http.Request, id string) {
	client, err := db.NewRedis()
	if err != nil {
		httpErr := errors.NewHTTPError(http.StatusInternalServerError, &errors.INTERNAL_SERVER_ERROR)
		httpErr.Write(w)
		return
	}

	defer client.Close()

	b, err := booleans.GetBoolean(client, r.Context(), id)
	if err != nil {
		httpErr, ok := err.(*errors.HTTPError)

		if !ok {
			httpErr = errors.NewHTTPError(http.StatusNotFound, &errors.NOT_FOUND_ERROR)
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

func handleToggleBooleanById(w http.ResponseWriter, r *http.Request, id string) {
	client, err := db.NewRedis()
	if err != nil {
		httpErr := errors.NewHTTPError(http.StatusInternalServerError, &errors.INTERNAL_SERVER_ERROR)
		httpErr.Write(w)
		return
	}

	defer client.Close()

	var b *booleans.Boolean
	if b, err = booleans.ToggleBoolean(client, r.Context(), id); err != nil {
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

func HandleBooleanById(w http.ResponseWriter, r *http.Request) {
	segments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	id := segments[len(segments)-1]

	if id == "" || id == "booleans" {
		httpErr := errors.NewHTTPError(http.StatusBadRequest, &errors.BAD_REQUEST_ERROR)
		httpErr.Write(w)
		return
	}

	params := r.URL.Query()
	params.Del("id")
	params.Add("id", id)

	r.URL.RawQuery = params.Encode()

	switch r.Method {
	case http.MethodGet:
		handleGetBooleanById(w, r, id)
	case http.MethodDelete:
		handleDeleteBooleanById(w, r, id)
	case http.MethodPatch:
		handleToggleBooleanById(w, r, id)
	case http.MethodPut:
		handleCreateBoolean(w, r)
	default:
		w.Header().Add("Allow", "GET, DELETE, PATCH, PUT")

		httpErr := errors.NewHTTPError(http.StatusMethodNotAllowed, &errors.METHOD_NOT_ALLOWED_ERROR)
		httpErr.Write(w)
	}
}
