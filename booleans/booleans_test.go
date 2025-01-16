package booleans

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/saschazar21/go-baas/db"
	"github.com/saschazar21/go-baas/test"
	"github.com/stretchr/testify/assert"
)

func TestValidateBooleanParams(t *testing.T) {
	type test struct {
		name    string
		data    BooleanParams
		wantErr bool
	}

	tests := []test{
		{
			name: "valid boolean params",
			data: BooleanParams{
				ExpiresAt: time.Now().Unix() + 20,
				ExpiresIn: 20,
				Id:        nil,
			},
			wantErr: false,
		},
		{
			name: "invalid expires_at",
			data: BooleanParams{
				ExpiresAt: time.Now().Unix() - 1,
				ExpiresIn: 20,
				Id:        nil,
			},
			wantErr: true,
		},
		{
			name: "invalid expires_in",
			data: BooleanParams{
				ExpiresAt: time.Now().Unix() + 20,
				ExpiresIn: -20,
				Id:        nil,
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.data.Validate(); (err != nil) != tc.wantErr {
				t.Errorf("BooleanParams.Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestValidateBoolean(t *testing.T) {
	type test struct {
		name    string
		data    Boolean
		wantErr bool
	}

	tests := []test{
		{
			name: "valid boolean",
			data: Boolean{
				Label: "test",
				Value: true,
				BooleanParams: &BooleanParams{
					ExpiresAt: time.Now().Unix() + 20,
					ExpiresIn: 20,
					Id:        nil,
				},
			},
			wantErr: false,
		},
		{
			name: "invalid expires_in",
			data: Boolean{
				Label: "test",
				BooleanParams: &BooleanParams{
					ExpiresAt: time.Now().Unix() + 20,
					ExpiresIn: -1,
					Id:        nil,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid expires_at",
			data: Boolean{
				Label: "test",
				Value: false,
				BooleanParams: &BooleanParams{
					ExpiresAt: time.Now().Unix() - 1,
					ExpiresIn: 20,
					Id:        nil,
				},
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.data.Validate(); (err != nil) != tc.wantErr {
				t.Errorf("Boolean.Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestParseBooleans(t *testing.T) {
	type testStruct struct {
		name    string
		headers http.Header
		params  url.Values
		body    []byte
		cmp     Boolean
		wantErr bool
	}

	tests := []testStruct{
		{
			name: "valid json boolean",
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			params: url.Values{
				"expires_at": []string{fmt.Sprintf("%d", time.Now().Unix()+20)},
				"expires_in": []string{"20"},
			},
			body: []byte(`{"label":"test","value":true}`),
			cmp: Boolean{
				Label: "test",
				Value: true,
				BooleanParams: &BooleanParams{
					ExpiresAt: time.Now().Unix() + 20,
					ExpiresIn: 20,
					Id:        nil,
				},
			},
			wantErr: false,
		},
		{
			name: "valid urlencoded boolean",
			headers: http.Header{
				"Content-Type": []string{"application/x-www-form-urlencoded"},
			},
			params: url.Values{
				"expires_at": []string{fmt.Sprintf("%d", time.Now().Unix()+20)},
				"expires_in": []string{"20"},
			},
			body: []byte(`label=test&value=true`),
			cmp: Boolean{
				Label: "test",
				Value: true,
				BooleanParams: &BooleanParams{
					ExpiresAt: time.Now().Unix() + 20,
					ExpiresIn: 20,
					Id:        nil,
				},
			},
			wantErr: false,
		},
		{
			name: "invalid content-type",
			headers: http.Header{
				"Content-Type": []string{"text/plain"},
			},
			params: url.Values{
				"expires_at": []string{fmt.Sprintf("%d", time.Now().Unix()+20)},
				"expires_in": []string{"20"},
			},
			body:    []byte(`{"label":"test","value":true}`),
			cmp:     Boolean{},
			wantErr: true,
		},
		{
			name: "invalid expires_in",
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			params: url.Values{
				"expires_at": []string{fmt.Sprintf("%d", time.Now().Unix()+20)},
				"expires_in": []string{"-20"},
			},
			body:    []byte(`{"label":"test","value":true}`),
			cmp:     Boolean{},
			wantErr: true,
		},
		{
			name: "invalid expires_at",
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			params: url.Values{
				"expires_at": []string{fmt.Sprintf("%d", time.Now().Unix()-1)},
				"expires_in": []string{"20"},
			},
			body:    []byte(`{"label":"test","value":true}`),
			cmp:     Boolean{},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buf := bytes.NewBuffer(tc.body)
			req := httptest.NewRequest("POST", "/api/v1/booleans", buf)
			req.Header = tc.headers
			req.URL.RawQuery = tc.params.Encode()

			b, err := ParseBoolean(req)

			if (err != nil) != tc.wantErr {
				t.Errorf("ParseBoolean() error = %v, wantErr %v", err, tc.wantErr)
			}

			if !tc.wantErr {
				assert.Equal(t, tc.cmp.Label, b.Label)
				assert.Equal(t, tc.cmp.Value, b.Value)
				assert.Equal(t, tc.cmp.ExpiresAt, b.ExpiresAt)
				assert.Equal(t, tc.cmp.ExpiresIn, b.ExpiresIn)
			}
		})
	}
}

func TestBooleans(t *testing.T) {
	inexistentId := "test_id"

	type testStruct struct {
		name    string
		data    Boolean
		wantErr bool
	}

	tests := []testStruct{
		{
			name: "valid boolean",
			data: Boolean{
				Label: "test",
				Value: true,
				BooleanParams: &BooleanParams{
					ExpiresAt: time.Now().Unix() + 20,
					ExpiresIn: 20,
					Id:        nil,
				},
			},
			wantErr: false,
		},
		{
			name: "invalid expires_in",
			data: Boolean{
				Label: "test",
				BooleanParams: &BooleanParams{
					ExpiresAt: time.Now().Unix() + 20,
					ExpiresIn: -1,
					Id:        nil,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid expires_at",
			data: Boolean{
				Label: "test",
				Value: false,
				BooleanParams: &BooleanParams{
					ExpiresAt: time.Now().Unix() - 1,
					ExpiresIn: 20,
					Id:        nil,
				},
			},
			wantErr: true,
		},
		{
			name: "inexistent boolean",
			data: Boolean{
				Label: "test",
				Value: false,
				BooleanParams: &BooleanParams{
					ExpiresAt: time.Now().Unix() + 20,
					ExpiresIn: 20,
					Id:        &inexistentId,
				},
			},
			wantErr: true,
		},
	}

	ctx := context.Background()

	container, err := test.CreateContainer(ctx, t)

	if err != nil {
		t.Fatalf("%v", err)
	}

	t.Cleanup(func() {
		test.TerminateContainer(container, t)
	})

	opts, err := redis.ParseURL(os.Getenv(db.REDIS_URL_ENV))

	if err != nil {
		t.Fatalf("%v", err)
	}

	rdb := redis.NewClient(opts)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.data.Save(rdb, ctx); (err != nil) != tc.wantErr {
				t.Errorf("Boolean.Save() error = %v, wantErr %v", err, tc.wantErr)
			}

			if !tc.wantErr {
				log.Println(tc.data)

				var b *Boolean
				if b, err = GetBoolean(rdb, ctx, *tc.data.Id); err != nil {
					t.Errorf("GetBoolean() error = %v", err)
				}

				assert.Equal(t, tc.data.Label, b.Label)
				assert.Equal(t, tc.data.Value, b.Value)

				log.Println(b)

				assert.Equal(t, int64(1), rdb.Exists(ctx, *tc.data.Id).Val())

				if b, err = ToggleBoolean(rdb, ctx, *tc.data.Id); err != nil {
					t.Errorf("ToggleBoolean() error = %v", err)
				}

				assert.Equal(t, !tc.data.Value, b.Value)

				if err := DeleteBoolean(rdb, ctx, *tc.data.Id); err != nil {
					t.Errorf("DeleteBoolean() error = %v", err)
				}

				if _, err = GetBoolean(rdb, ctx, *tc.data.Id); err == nil {
					t.Errorf("GetBoolean() error = %v, wantErr %v", err, true)
				}
			}
		})
	}

}
