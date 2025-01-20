package v1_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	v1 "github.com/saschazar21/go-baas/api/v1"
	"github.com/saschazar21/go-baas/booleans"
	"github.com/saschazar21/go-baas/errors"
	"github.com/saschazar21/go-baas/test"
	"github.com/stretchr/testify/assert"
)

type booleanResponse struct {
	Data struct {
		Id string `json:"id"`
		*booleans.Boolean
	} `json:"data"`
}

func TestHandleBooleans(t *testing.T) {
	type testStruct struct {
		name       string
		method     string
		parameters url.Values
		data       booleans.Boolean
		wantErr    bool
	}

	ctx := context.Background()

	container, err := test.CreateContainer(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		test.TerminateContainer(container, t)
	})

	tests := []testStruct{
		{
			name:   "valid boolean",
			method: http.MethodPost,
			parameters: url.Values{
				"expires_at": []string{fmt.Sprintf("%d", time.Now().Unix()+20)},
				"expires_in": []string{fmt.Sprintf("%d", 20)},
				"id":         nil,
			},
			data: booleans.Boolean{
				Label: "test",
				Value: true,
			},
			wantErr: false,
		},
		{
			name:   "valid boolean without label",
			method: http.MethodPost,
			parameters: url.Values{
				"expires_at": []string{fmt.Sprintf("%d", time.Now().Unix()+20)},
				"expires_in": []string{fmt.Sprintf("%d", 20)},
				"id":         nil,
			},
			data: booleans.Boolean{
				Value: true,
			},
			wantErr: false,
		},
		{
			name:   "valid boolean without value",
			method: http.MethodPost,
			parameters: url.Values{
				"expires_at": []string{fmt.Sprintf("%d", time.Now().Unix()+20)},
				"expires_in": []string{fmt.Sprintf("%d", 20)},
				"id":         nil,
			},
			data:    booleans.Boolean{},
			wantErr: false,
		},
		{
			name:   "invalid expires_in",
			method: http.MethodPost,
			parameters: url.Values{
				"expires_at": []string{fmt.Sprintf("%d", time.Now().Unix()+20)},
				"expires_in": []string{fmt.Sprintf("%d", -1)},
				"id":         nil,
			},
			data: booleans.Boolean{
				Label: "test",
				Value: true,
			},
			wantErr: true,
		},
		{
			name:   "invalid expires_at",
			method: http.MethodPost,
			parameters: url.Values{
				"expires_at": []string{fmt.Sprintf("%d", time.Now().Unix()-1)},
				"expires_in": []string{fmt.Sprintf("%d", 20)},
				"id":         nil,
			},
			data: booleans.Boolean{
				Label: "test",
				Value: true,
			},
			wantErr: true,
		},
		{
			name:       "invalid method",
			method:     http.MethodGet,
			parameters: url.Values{},
			data:       booleans.Boolean{},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			var server *httptest.Server

			t.Cleanup(func() {
				server.Close()
			})

			server = httptest.NewServer(http.HandlerFunc(v1.HandleBooleans))

			var u *url.URL
			if u, err = url.Parse(server.URL); err != nil {
				t.Fatal(err)
			}

			u.RawQuery = tt.parameters.Encode()

			var body []byte
			if body, err = json.Marshal(tt.data); err != nil {
				t.Fatal(err)
			}

			var req *http.Request
			if req, err = http.NewRequest(tt.method, u.String(), bytes.NewBuffer(body)); err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")

			client := server.Client()

			var resp *http.Response
			if resp, err = client.Do(req); err != nil {
				t.Fatal(err)
			}

			if !tt.wantErr {
				assert.Equal(t, http.StatusOK, resp.StatusCode)

				var res booleanResponse
				if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, tt.data.Label, res.Data.Label)
				assert.Equal(t, tt.data.Value, res.Data.Value)
				assert.NotEmpty(t, res.Data.Id)
			} else {
				assert.NotEqual(t, http.StatusOK, resp.StatusCode)

				var httpErr errors.HTTPError
				if err = json.NewDecoder(resp.Body).Decode(&httpErr); err != nil {
					t.Fatal(err)
				}

				assert.NotEmpty(t, httpErr.Errors)
				assert.Equal(t, (*httpErr.Errors)[0].Status, resp.StatusCode)
			}
		})
	}
}
