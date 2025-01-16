package booleans

import (
	"testing"
	"time"
)

func TestCustomValidator(t *testing.T) {
	type test struct {
		name    string
		data    interface{}
		wantErr bool
	}

	tests := []test{
		{
			name: "epoch is greater than now",
			data: struct {
				Val int64 `validate:"epoch-gt-now"`
			}{time.Now().Unix() + 1},
			wantErr: false,
		},
		{
			name: "epoch is equal to now",
			data: struct {
				Val int64 `validate:"epoch-gt-now"`
			}{time.Now().Unix()},
			wantErr: true,
		},
		{
			name: "epoch is less than now",
			data: struct {
				Val int64 `validate:"epoch-gt-now"`
			}{time.Now().Unix() - 1},
			wantErr: true,
		},
		{
			name: "epoch is not an integer",
			data: struct {
				Val string `validate:"epoch-gt-now"`
			}{"not an integer"},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := CustomValidateStruct(tc.data); (err != nil) != tc.wantErr {
				t.Errorf("CustomValidateStruct() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
