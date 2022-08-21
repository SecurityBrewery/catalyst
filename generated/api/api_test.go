package api

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func Test_parseQueryOptionalBoolArray(t *testing.T) {
	type args struct {
		r   *http.Request
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    []bool
		wantErr bool
	}{
		{
			name: "bool array",
			args: args{
				r: httptest.NewRequest(
					http.MethodGet,
					"https://try.catalyst-soar.com/api/tickets?type=alert&offset=0&count=10&sort=status%2Cowner%2Ccreated&desc=true%2Cfalse%2Cfalse&query=status+%3D%3D+%27open%27+AND+%28owner+%3D%3D+%27eve%27+OR+%21owner%29",
					nil,
				),
				key: "desc",
			},
			want:    []bool{true, false, false},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseQueryOptionalBoolArray(tt.args.r, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseQueryOptionalBoolArray() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseQueryOptionalBoolArray() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseQueryOptionalStringArray(t *testing.T) {
	type args struct {
		r   *http.Request
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "string array",
			args: args{
				r: httptest.NewRequest(
					http.MethodGet,
					"https://try.catalyst-soar.com/api/tickets?type=alert&offset=0&count=10&sort=status%2Cowner%2Ccreated&desc=true%2Cfalse%2Cfalse&query=status+%3D%3D+%27open%27+AND+%28owner+%3D%3D+%27eve%27+OR+%21owner%29",
					nil,
				),
				key: "sort",
			},
			want:    []string{"status", "owner", "created"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseQueryOptionalStringArray(tt.args.r, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseQueryOptionalStringArray() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseQueryOptionalStringArray() got = %v, want %v", got, tt.want)
			}
		})
	}
}