package forwarded

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		forwarded string
		want      Forwarded
		wantErr   bool
	}{
		{
			forwarded: "",
			want:      Forwarded{},
		},
		{
			forwarded: "for=192.0.2.43, for=198.51.100.17",
			want: Forwarded{
				{
					For: "192.0.2.43",
				},
				{
					For: "198.51.100.17",
				},
			},
		},
		{
			forwarded: "for=192.0.2.60;proto=http ; by = 203.0.113.43",
			want: Forwarded{
				{
					For:   "192.0.2.60",
					Proto: "http",
					By:    "203.0.113.43",
				},
			},
		},
		{
			forwarded: "for = _hidden, for=_SEVKISEK",
			want: Forwarded{
				{
					For: "_hidden",
				},
				{
					For: "_SEVKISEK",
				},
			},
		},
		{
			forwarded: "for=192.0.2.43,for=198.51.100.17;by=203.0.113.60;proto=http;host=example.com",
			want: Forwarded{
				{
					For: "192.0.2.43",
				},
				{
					For:   "198.51.100.17",
					By:    "203.0.113.60",
					Proto: "http",
					Host:  "example.com",
				},
			},
		},
		{
			forwarded: "for=192.0.2.43, for=\"[2001:db8:cafe::17]\", for=unknown",
			want: Forwarded{
				{
					For: "192.0.2.43",
				},
				{
					For: "[2001:db8:cafe::17]",
				},
				{
					For: "unknown",
				},
			},
		},
		{
			forwarded: "For=\"[2001:db8:cafe::17]:4711\"",
			want: Forwarded{
				{
					For: "[2001:db8:cafe::17]:4711",
				},
			},
		},
		{
			forwarded: "for=\"_gazonk\"",
			want: Forwarded{
				{
					For: "_gazonk",
				},
			},
		},
		{
			forwarded: "for=\"\"",
			want: Forwarded{
				{},
			},
		},
		{
			forwarded: "by=\"\\\"\"",
			want: Forwarded{
				{
					By: "\"",
				},
			},
		},
		{
			forwarded: "for=",
			want:      Forwarded{},
			wantErr:   true,
		},
		{
			forwarded: "bad=xxx",
			want:      Forwarded{},
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.forwarded, func(t *testing.T) {
			got, err := Parse(tt.forwarded)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func Test_parserError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "test",
			err:  newParserError(0, "test"),
			want: "1: test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("parserError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
