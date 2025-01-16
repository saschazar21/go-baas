package booleans

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"

	"github.com/saschazar21/go-baas/errors"
)

const base58 = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

func generateRandomId() string {
	data := make([]byte, 16)

	for i := range data {
		data[i] = base58[rand.Int63()%int64(len(base58))]
	}

	return string(data)
}

func parseJsonEncodedBody(r *http.Request, d interface{}) (err error) {
	if r.Header.Get("Content-Type") != "application/json" {
		return errors.NewHTTPError(http.StatusUnsupportedMediaType, &errors.UNSUPPORTED_MEDIA_TYPE_ERROR)
	}

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(d); err != nil {
		log.Println(err)

		return errors.NewHTTPError(http.StatusBadRequest, &errors.BAD_REQUEST_ERROR)
	}

	return
}

func parseUrlEncodedBody(r *http.Request, d interface{}) (err error) {
	if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		return errors.NewHTTPError(http.StatusUnsupportedMediaType, &errors.UNSUPPORTED_MEDIA_TYPE_ERROR)
	}

	if err := r.ParseForm(); err != nil {
		log.Println(err)

		return errors.NewHTTPError(http.StatusBadRequest, &errors.BAD_REQUEST_ERROR)
	}

	if err := decoder.Decode(d, r.PostForm); err != nil {
		log.Println(err)

		return errors.NewHTTPError(http.StatusBadRequest, &errors.BAD_REQUEST_ERROR)
	}

	return
}
