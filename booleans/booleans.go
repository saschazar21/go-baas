package booleans

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/saschazar21/go-baas/errors"
)

const (
	BOOLEAN_LABEL = "label"
	BOOLEAN_VALUE = "value"
)

type booleanWithId struct {
	Id string `json:"id"`
	*Boolean
}

type booleanResponse struct {
	Data *booleanWithId `json:"data"`
}

type BooleanParams struct {
	ExpiresAt int64   `schema:"expires_at" validate:"omitempty,epoch-gt-now"`
	ExpiresIn int64   `schema:"expires_in" validate:"omitempty,gt=0"`
	Id        *string `schema:"id"`
}

func (b *BooleanParams) Validate() (err error) {
	if err = CustomValidateStruct(b); err != nil {
		log.Println(err)

		return errors.NewHTTPError(http.StatusBadRequest, &errors.BAD_REQUEST_ERROR)
	}

	return
}

type Boolean struct {
	Label string `json:"label,omitempty" redis:"label" schema:"label"`
	Value bool   `json:"value" redis:"value" schema:"value"`

	*BooleanParams `json:"-" redis:"-" schema:"-" validate:"omitempty"`
}

func (b Boolean) String() string {
	var id *string

	if b.BooleanParams != nil && b.Id != nil {
		id = b.Id
	} else {
		i := "<nil>"
		id = &i
	}

	return fmt.Sprintf("[%s]: %t, Label: \"%s\"", *id, b.Value, b.Label)
}

func (b *Boolean) Save(client *redis.Client, ctx context.Context) (err error) {
	if err = b.Validate(); err != nil {
		return
	}

	if b.BooleanParams == nil {
		b.BooleanParams = &BooleanParams{}
	}

	if b.Id != nil && client.Exists(ctx, *b.Id).Val() == 0 {
		log.Printf("Boolean with ID %s not found", *b.Id)

		return errors.NewHTTPError(http.StatusNotFound, &errors.NOT_FOUND_ERROR)
	}

	if b.Id == nil {
		id := generateRandomId()

		for client.Exists(ctx, id).Val() > 0 {
			id = generateRandomId()
		}

		b.Id = &id
	}

	if err = client.HSet(ctx, *b.Id, b).Err(); err != nil {
		log.Println(err)

		return errors.NewHTTPError(http.StatusInternalServerError, &errors.INTERNAL_SERVER_ERROR)
	}

	if b.ExpiresIn > 0 || b.ExpiresAt > 0 {
		var ttl int64

		if b.ExpiresAt > 0 {
			ttl = b.ExpiresAt
		} else {
			ttl = b.ExpiresIn + time.Now().Unix()
		}

		if err = client.ExpireAt(ctx, *b.Id, time.Unix(ttl, 0)).Err(); err != nil {
			log.Println(err)

			return errors.NewHTTPError(http.StatusInternalServerError, &errors.INTERNAL_SERVER_ERROR)
		}
	}

	return
}

func (b *Boolean) Validate() (err error) {
	if err = CustomValidateStruct(b); err != nil {
		log.Println(err)

		return errors.NewHTTPError(http.StatusBadRequest, &errors.BAD_REQUEST_ERROR)
	}

	return
}

func DeleteBoolean(client *redis.Client, ctx context.Context, id string) (err error) {
	if err = client.Del(ctx, id).Err(); err != nil {
		log.Println(err)

		return errors.NewHTTPError(http.StatusInternalServerError, &errors.INTERNAL_SERVER_ERROR)
	}

	return
}

func GetBoolean(client *redis.Client, ctx context.Context, id string) (b *Boolean, err error) {
	if client.Exists(ctx, id).Val() == 0 {
		log.Printf("Boolean with ID %s not found", id)

		return b, errors.NewHTTPError(http.StatusNotFound, &errors.NOT_FOUND_ERROR)
	}

	b = new(Boolean)

	if err = client.HGetAll(ctx, id).Scan(b); err != nil {
		log.Println(err)

		return b, errors.NewHTTPError(http.StatusInternalServerError, &errors.INTERNAL_SERVER_ERROR)
	}

	b.BooleanParams = &BooleanParams{
		Id: &id,
	}

	return
}

func ToggleBoolean(client *redis.Client, ctx context.Context, id string) (b *Boolean, err error) {
	if b, err = GetBoolean(client, ctx, id); err != nil {
		return
	}

	b.Value = !b.Value

	if err = client.HSet(ctx, id, BOOLEAN_VALUE, b.Value).Err(); err != nil {
		log.Println(err)

		return nil, errors.NewHTTPError(http.StatusInternalServerError, &errors.INTERNAL_SERVER_ERROR)
	}

	return
}

func ParseBoolean(r *http.Request) (b *Boolean, err error) {
	var params BooleanParams

	if err = decoder.Decode(&params, r.URL.Query()); err != nil {
		log.Println(err)

		return b, errors.NewHTTPError(http.StatusBadRequest, &errors.BAD_REQUEST_ERROR)
	}

	if err = params.Validate(); err != nil {
		return
	}

	b = new(Boolean)

	switch r.Header.Get("Content-Type") {
	case "application/json":
		if err = parseJsonEncodedBody(r, b); err != nil {
			return
		}
	case "application/x-www-form-urlencoded":
		if err = parseUrlEncodedBody(r, b); err != nil {
			return
		}
	default:
		err = errors.NewHTTPError(http.StatusUnsupportedMediaType, &errors.UNSUPPORTED_MEDIA_TYPE_ERROR)
		return
	}

	b.BooleanParams = &params

	if err = b.Validate(); err != nil {
		return
	}

	return
}

func CreateBooleanResponse(b *Boolean) (body *booleanResponse) {
	body = &booleanResponse{
		Data: &booleanWithId{
			Id:      *b.Id,
			Boolean: b,
		},
	}

	return
}
