package transfer_money

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func Test_validate(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    TransferDTO
		wantErr bool
	}{
		{
			name: "err on empty request",
			args: args{
				r: httptest.NewRequest("", "/", strings.NewReader("")),
			},
			want:    TransferDTO{},
			wantErr: true,
		},
		{
			name: "err on empty from_wallet_id",
			args: args{
				r: httptest.NewRequest("", "/", strings.NewReader("{}")),
			},
			want:    TransferDTO{},
			wantErr: true,
		},
		{
			name: "err on empty to_wallet_id",
			args: args{
				r: httptest.NewRequest("", "/", strings.NewReader("{\"from_wallet_id\": 1}")),
			},
			want:    TransferDTO{},
			wantErr: true,
		},
		{
			name: "err on empty value",
			args: args{
				r: httptest.NewRequest("", "/", strings.NewReader("{\"from_wallet_id\": 1, \"to_wallet_id\": 2}")),
			},
			want:    TransferDTO{},
			wantErr: true,
		},
		{
			name: "err on empty idempotency_key",
			args: args{
				r: httptest.NewRequest("", "/", strings.NewReader("{\"from_wallet_id\": 1, \"to_wallet_id\": 2, \"value\": 100.53}")),
			},
			want:    TransferDTO{},
			wantErr: true,
		},
		{
			name: "no err",
			args: args{
				r: httptest.NewRequest("", "/", strings.NewReader("{\"from_wallet_id\": 1, \"to_wallet_id\": 2, \"value\": 100.53, \"idempotency_key\": \"boo\"}")),
			},
			want: TransferDTO{
				FromWalletID:   1,
				ToWalletID:     2,
				Value:          100.53,
				IdempotencyKey: "boo",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validate(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validate() got = %v, want %v", got, tt.want)
			}
		})
	}
}
