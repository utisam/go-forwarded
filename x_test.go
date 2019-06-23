package forwarded

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFromX(t *testing.T) {
	tests := []struct {
		name    string
		args    []XField
		want    Forwarded
		wantErr bool
	}{
		{
			args: []XField{
				ForwardedFor("203.0.113.1, 203.0.113.2"),
			},
			want: Forwarded{
				{For: "203.0.113.1"},
				{For: "203.0.113.2"},
			},
		},
		{
			name: "empty",
			want: Forwarded{},
		},
		{
			name: "invalid length",
			args: []XField{
				ForwardedFor("203.0.113.1, 203.0.113.2"),
				ForwardedHost("example.com"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromX(tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromX() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromX() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleFromX() {
	forwarded, _ := FromX(
		ForwardedFor("203.0.113.1, 203.0.113.2"),
		ForwardedHost("example.com, example.org"),
	)
	fmt.Println("Forwarded: " + forwarded.String())
	// Output:
	// Forwarded: for=203.0.113.1;host=example.com,for=203.0.113.2;host=example.org
}

func TestAlignX(t *testing.T) {
	tests := []struct {
		name string
		args []XField
		want Forwarded
	}{
		{
			args: []XField{
				ForwardedFor("203.0.113.1, 203.0.113.2"),
			},
			want: Forwarded{
				{For: "203.0.113.1"},
				{For: "203.0.113.2"},
			},
		},
		{
			name: "empty",
			want: Forwarded{},
		},
		{
			name: "overwrite",
			args: []XField{
				ForwardedFor("203.0.113.1, 203.0.113.2"),
				RealIP("198.51.100.1"),
			},
			want: Forwarded{
				{For: "198.51.100.1"},
				{For: "203.0.113.2"},
			},
		},
		{
			name: "different length",
			args: []XField{
				ForwardedHost("example.com"),
				ForwardedFor("203.0.113.1, 203.0.113.2"),
			},
			want: Forwarded{
				{For: "203.0.113.1", Host: "example.com"},
				{For: "203.0.113.2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AlignX(tt.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AlignX() = %v, want %v", got, tt.want)
			}
		})
	}
}
