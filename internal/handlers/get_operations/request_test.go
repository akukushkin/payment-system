package get_operations

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
		want    FilterDTO
		wantErr bool
	}{
		{
			name: "err on empty request",
			args: args{
				r: httptest.NewRequest("", "/", strings.NewReader("")),
			},
			want:    FilterDTO{},
			wantErr: true,
		},
		{
			name: "err on empty wallet_id",
			args: args{
				r: httptest.NewRequest("", "/", strings.NewReader("{}")),
			},
			want:    FilterDTO{},
			wantErr: true,
		},
		{
			name: "err on empty date",
			args: args{
				r: httptest.NewRequest("", "/", strings.NewReader("{\"wallet_id\": 1}")),
			},
			want:    FilterDTO{},
			wantErr: true,
		},
		{
			name: "err on invalid date",
			args: args{
				r: httptest.NewRequest("", "/", strings.NewReader("{\"wallet_id\": 1, \"date\": \"boo\"}")),
			},
			want:    FilterDTO{},
			wantErr: true,
		},
		{
			name: "no err",
			args: args{
				r: httptest.NewRequest("", "/", strings.NewReader("{\"wallet_id\": 1, \"date\": \"2021-06-30\", \"direction\":1}")),
			},
			want: FilterDTO{
				WalletID:  1,
				Date:      "2021-06-30",
				Direction: 1,
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
