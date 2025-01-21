package v1_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	rdb "github.com/redis/go-redis/v9"
	v1 "github.com/saschazar21/go-baas/api/v1"
	"github.com/saschazar21/go-baas/booleans"
	"github.com/saschazar21/go-baas/db"
	"github.com/saschazar21/go-baas/errors"
	"github.com/saschazar21/go-baas/test"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

const (
	BOOLEAN_TEST_ID = "test"
)

func TestHandleBooleanById(t *testing.T) {
	var client *rdb.Client
	var container *redis.RedisContainer
	var err error
	var server *httptest.Server

	ctx := context.Background()

	t.Cleanup(func() {
		client.Close()
		server.Close()
		test.TerminateContainer(container, t)
	})

	if container, err = test.CreateContainer(ctx, t); err != nil {
		t.Fatal(err)
	}

	server = httptest.NewServer(http.HandlerFunc(v1.HandleBooleanById))
	if client, err = db.NewRedis(); err != nil {
		t.Fatal(err)
	}

	t.Run("get boolean by id", func(t *testing.T) {
		t.Cleanup(func() {
			if err = client.Del(ctx, BOOLEAN_TEST_ID).Err(); err != nil {
				t.Fatal(err)
			}
		})

		if err = client.HSet(ctx, BOOLEAN_TEST_ID, &booleans.Boolean{
			Label: BOOLEAN_TEST_ID,
			Value: true,
		}).Err(); err != nil {
			t.Fatal(err)
		}

		tests := []struct {
			name    string
			method  string
			id      string
			wantErr bool
		}{
			{
				name:    "get boolean by id",
				method:  http.MethodGet,
				id:      BOOLEAN_TEST_ID,
				wantErr: false,
			},
			{
				name:    "get inexistent boolean by id",
				method:  http.MethodGet,
				id:      "inexistentId",
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req, err := http.NewRequest(tt.method, server.URL, nil)
				if err != nil {
					t.Fatal(err)
				}

				req.URL.Path = "/api/v1/boolean/" + tt.id
				req.URL.RawQuery = "id=" + tt.id

				httpClient := server.Client()
				res, err := httpClient.Do(req)
				if err != nil {
					t.Fatal(err)
				}

				if !tt.wantErr {
					assert.Equal(t, http.StatusOK, res.StatusCode)

					var b booleanResponse
					if err = json.NewDecoder(res.Body).Decode(&b); err != nil {
						t.Fatal(err)
					}

					assert.Equal(t, BOOLEAN_TEST_ID, b.Data.Label)
					assert.Equal(t, true, b.Data.Value)
				} else {
					assert.Equal(t, http.StatusNotFound, res.StatusCode)

					var httpErr errors.HTTPError
					if err = json.NewDecoder(res.Body).Decode(&httpErr); err != nil {
						t.Fatal(err)
					}

					assert.NotEmpty(t, httpErr.Errors)
					assert.Equal(t, (*httpErr.Errors)[0].Status, res.StatusCode)
				}
			})
		}
	})

	t.Run("delete boolean by id", func(t *testing.T) {
		if err = client.HSet(ctx, BOOLEAN_TEST_ID, &booleans.Boolean{
			Label: BOOLEAN_TEST_ID,
			Value: true,
		}).Err(); err != nil {
			t.Fatal(err)
		}

		tests := []struct {
			name    string
			method  string
			id      string
			wantErr bool
		}{
			{
				name:    "delete boolean by id",
				method:  http.MethodDelete,
				id:      BOOLEAN_TEST_ID,
				wantErr: false,
			},
			{
				name:    "delete inexistent boolean by id",
				method:  http.MethodDelete,
				id:      "inexistentId",
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req, err := http.NewRequest(tt.method, server.URL, nil)
				if err != nil {
					t.Fatal(err)
				}

				req.URL.Path = "/api/v1/boolean/" + tt.id
				req.URL.RawQuery = "id=" + tt.id

				httpClient := server.Client()
				res, err := httpClient.Do(req)
				if err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, http.StatusNoContent, res.StatusCode)
			})
		}
	})

	t.Run("toggle boolean by id", func(t *testing.T) {
		t.Cleanup(func() {
			if err = client.Del(ctx, BOOLEAN_TEST_ID).Err(); err != nil {
				t.Fatal(err)
			}
		})

		if err = client.HSet(ctx, BOOLEAN_TEST_ID, &booleans.Boolean{
			Label: BOOLEAN_TEST_ID,
			Value: true,
		}).Err(); err != nil {
			t.Fatal(err)
		}

		tests := []struct {
			name    string
			method  string
			id      string
			wantErr bool
		}{
			{
				name:    "toggle boolean by id",
				method:  http.MethodPatch,
				id:      BOOLEAN_TEST_ID,
				wantErr: false,
			},
			{
				name:    "toggle inexistent boolean by id",
				method:  http.MethodPatch,
				id:      "inexistentId",
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req, err := http.NewRequest(tt.method, server.URL, nil)
				if err != nil {
					t.Fatal(err)
				}

				req.URL.Path = "/api/v1/boolean/" + tt.id
				req.URL.RawQuery = "id=" + tt.id

				httpClient := server.Client()
				res, err := httpClient.Do(req)
				if err != nil {
					t.Fatal(err)
				}

				if !tt.wantErr {
					assert.Equal(t, http.StatusOK, res.StatusCode)

					var b booleanResponse
					if err = json.NewDecoder(res.Body).Decode(&b); err != nil {
						t.Fatal(err)
					}

					assert.Equal(t, BOOLEAN_TEST_ID, b.Data.Label)
					assert.Equal(t, false, b.Data.Value)
				} else {
					assert.Equal(t, http.StatusNotFound, res.StatusCode)

					var httpErr errors.HTTPError
					if err = json.NewDecoder(res.Body).Decode(&httpErr); err != nil {
						t.Fatal(err)
					}

					assert.NotEmpty(t, httpErr.Errors)
					assert.Equal(t, (*httpErr.Errors)[0].Status, res.StatusCode)
				}
			})
		}
	})

	t.Run("update boolean by id", func(t *testing.T) {
		t.Cleanup(func() {
			if err = client.Del(ctx, BOOLEAN_TEST_ID).Err(); err != nil {
				t.Fatal(err)
			}
		})

		if err = client.HSet(ctx, BOOLEAN_TEST_ID, &booleans.Boolean{
			Label: BOOLEAN_TEST_ID,
			Value: true,
		}).Err(); err != nil {
			t.Fatal(err)
		}

		tests := []struct {
			name    string
			method  string
			id      string
			data    booleans.Boolean
			wantErr bool
		}{
			{
				name:   "update boolean by id",
				method: http.MethodPut,
				id:     BOOLEAN_TEST_ID,
				data: booleans.Boolean{
					Label: fmt.Sprintf("%s-updated", BOOLEAN_TEST_ID),
					Value: false,
				},
				wantErr: false,
			},
			{
				name:   "update inexistent boolean by id",
				method: http.MethodPut,
				id:     "inexistentId",
				data: booleans.Boolean{
					Label: "inexistentId",
					Value: false,
				},
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				encoded, err := json.Marshal(tt.data)
				if err != nil {
					t.Fatal(err)
				}

				req, err := http.NewRequest(tt.method, server.URL, bytes.NewBuffer(encoded))
				if err != nil {
					t.Fatal(err)
				}

				req.URL.Path = "/api/v1/boolean/" + tt.id
				req.URL.RawQuery = "id=" + tt.id

				req.Header.Add("Content-Type", "application/json")

				httpClient := server.Client()
				res, err := httpClient.Do(req)
				if err != nil {
					t.Fatal(err)
				}

				if !tt.wantErr {
					assert.Equal(t, http.StatusOK, res.StatusCode)

					var b booleanResponse
					if err = json.NewDecoder(res.Body).Decode(&b); err != nil {
						t.Fatal(err)
					}

					assert.Equal(t, tt.data.Label, b.Data.Label)
					assert.Equal(t, tt.data.Value, b.Data.Value)
				} else {
					assert.Equal(t, http.StatusNotFound, res.StatusCode)

					var httpErr errors.HTTPError
					if err = json.NewDecoder(res.Body).Decode(&httpErr); err != nil {
						t.Fatal(err)
					}

					assert.NotEmpty(t, httpErr.Errors)
					assert.Equal(t, (*httpErr.Errors)[0].Status, res.StatusCode)
				}
			})
		}
	})

	t.Run("global request error handling", func(t *testing.T) {
		tests := []struct {
			name    string
			method  string
			id      string
			wantErr int
		}{
			{
				name:    "unsupported method",
				method:  http.MethodPost,
				id:      BOOLEAN_TEST_ID,
				wantErr: http.StatusMethodNotAllowed,
			},
			{
				name:    "missing id",
				method:  http.MethodGet,
				id:      "",
				wantErr: http.StatusNotFound,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req, err := http.NewRequest(tt.method, server.URL, nil)
				if err != nil {
					t.Fatal(err)
				}

				req.URL.Path = "/api/v1/boolean/" + tt.id

				if tt.id != "" {
					req.URL.RawQuery = "id=" + tt.id
				}

				httpClient := server.Client()
				res, err := httpClient.Do(req)
				if err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, tt.wantErr, res.StatusCode)

				var httpErr errors.HTTPError
				if err = json.NewDecoder(res.Body).Decode(&httpErr); err != nil {
					t.Fatal(err)
				}

				assert.NotEmpty(t, httpErr.Errors)
				assert.Equal(t, (*httpErr.Errors)[0].Status, res.StatusCode)
			})
		}
	})
}
