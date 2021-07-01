package deposit_money

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
		want    DepositDTO
		wantErr bool
	}{
		{
			name: "err on empty request",
			args: args{
				r: httptest.NewRequest("", "/", strings.NewReader("")),
			},
			want:    DepositDTO{},
			wantErr: true,
		},
		{
			name: "err on empty idempotency_key",
			args: args{
				r: httptest.NewRequest("", "/", strings.NewReader("{}")),
			},
			want:    DepositDTO{},
			wantErr: true,
		},
		{
			name: "err on empty wallet_id",
			args: args{
				r: httptest.NewRequest("", "/", strings.NewReader("{\"idempotency_key\": \"foo\"}")),
			},
			want:    DepositDTO{},
			wantErr: true,
		},
		{
			name: "err on empty value",
			args: args{
				r: httptest.NewRequest("", "/", strings.NewReader("{\"idempotency_key\": \"foo\", \"wallet_id\": 1}")),
			},
			want:    DepositDTO{},
			wantErr: true,
		},
		{
			name: "no err",
			args: args{
				r: httptest.NewRequest("", "/", strings.NewReader("{\"idempotency_key\": \"foo\", \"wallet_id\": 1, \"value\": 100.55}")),
			},
			want: DepositDTO{
				IdempotencyKey: "foo",
				WalletID:       1,
				Value:          100.55,
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
