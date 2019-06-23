package forwarded_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/utisam/go-forwarded"
)

func TestFromX(t *testing.T) {
	tests := []struct {
		name    string
		args    []forwarded.XField
		want    forwarded.Forwarded
		wantErr bool
	}{
		{
			args: []forwarded.XField{
				forwarded.For("203.0.113.1, 203.0.113.2"),
			},
			want: forwarded.Forwarded{
				{For: "203.0.113.1"},
				{For: "203.0.113.2"},
			},
		},
		{
			name: "empty",
			want: forwarded.Forwarded{},
		},
		{
			name: "invalid length",
			args: []forwarded.XField{
				forwarded.For("203.0.113.1, 203.0.113.2"),
				forwarded.Host("example.com"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := forwarded.FromX(tt.args...)
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
	f, _ := forwarded.FromX(
		forwarded.For("203.0.113.1, 203.0.113.2"),
		forwarded.Host("example.com, example.org"),
	)
	fmt.Println("Forwarded: " + f.String())
	// Output:
	// Forwarded: for=203.0.113.1;host=example.com,for=203.0.113.2;host=example.org
}

func TestAlignX(t *testing.T) {
	tests := []struct {
		name string
		args []forwarded.XField
		want forwarded.Forwarded
	}{
		{
			args: []forwarded.XField{
				forwarded.For("203.0.113.1, 203.0.113.2"),
			},
			want: forwarded.Forwarded{
				{For: "203.0.113.1"},
				{For: "203.0.113.2"},
			},
		},
		{
			name: "empty",
			want: forwarded.Forwarded{},
		},
		{
			name: "overwrite",
			args: []forwarded.XField{
				forwarded.For("203.0.113.1, 203.0.113.2"),
				forwarded.RealIP("198.51.100.1"),
			},
			want: forwarded.Forwarded{
				{For: "198.51.100.1"},
				{For: "203.0.113.2"},
			},
		},
		{
			name: "different length",
			args: []forwarded.XField{
				forwarded.Host("example.com"),
				forwarded.For("203.0.113.1, 203.0.113.2"),
			},
			want: forwarded.Forwarded{
				{For: "203.0.113.1", Host: "example.com"},
				{For: "203.0.113.2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := forwarded.AlignX(tt.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AlignX() = %v, want %v", got, tt.want)
			}
		})
	}
}
