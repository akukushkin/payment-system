package add_wallet

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
		want    WalletInDTO
		wantErr bool
	}{
		{
			name: "err on empty request body",
			args: args{
				r: httptest.NewRequest("", "/", strings.NewReader("")),
			},
			want:    WalletInDTO{},
			wantErr: true,
		},
		{
			name: "err on empty idempotency_key",
			args: args{
				r: httptest.NewRequest("", "/", strings.NewReader("{}")),
			},
			want:    WalletInDTO{},
			wantErr: true,
		},
		{
			name: "no err",
			args: args{
				r: httptest.NewRequest("", "/", strings.NewReader("{\"idempotency_key\": \"foo\"}")),
			},
			want: WalletInDTO{
				IdempotencyKey: "foo",
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
