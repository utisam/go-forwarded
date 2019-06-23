package forwarded_test

import (
	"fmt"
	"testing"

	"github.com/utisam/go-forwarded"
)

func Example() {
	f, _ := forwarded.Parse("for=192.0.2.43, for=\"[2001:db8:cafe::17]\", for=unknown")
	fmt.Printf("%s\n", f)
	fmt.Printf("% s\n", f) // With spaces
	// Output:
	// for=192.0.2.43,for=[2001:db8:cafe::17],for=unknown
	// for=192.0.2.43, for=[2001:db8:cafe::17], for=unknown
}

func TestForwarded_String(t *testing.T) {
	tests := []struct {
		name      string
		forwarded forwarded.Forwarded
		want      string
	}{
		{
			name: "all",
			forwarded: forwarded.Forwarded{{
				By:    "192.168.1.1",
				For:   "203.0.113.1",
				Host:  "example.com",
				Proto: "http",
			}},
			want: "by=192.168.1.1;for=203.0.113.1;host=example.com;proto=http",
		},
		{
			name: "empty",
			want: "",
		},
		{
			name: "for",
			forwarded: forwarded.Forwarded{{
				For: "203.0.113.1",
			}},
			want: "for=203.0.113.1",
		},
		{
			name: "include space",
			forwarded: forwarded.Forwarded{{
				Host: "l7 loadbalancer",
			}},
			want: "host=\"l7 loadbalancer\"",
		},
		{
			name: "include double quote",
			forwarded: forwarded.Forwarded{{
				Host: "the \"test\" server",
			}},
			want: "host=\"the \\\"test\\\" server\"",
		},
		{
			name: "2 for",
			forwarded: forwarded.Forwarded{
				{For: "203.0.113.1"},
				{For: "203.0.113.2"},
			},
			want: "for=203.0.113.1,for=203.0.113.2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.forwarded.String(); got != tt.want {
				t.Errorf("Forwarded.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestForwarded_StringSpace(t *testing.T) {
	tests := []struct {
		name      string
		forwarded forwarded.Forwarded
		want      string
	}{
		{
			name: "all",
			forwarded: forwarded.Forwarded{{
				By:    "192.168.1.1",
				For:   "203.0.113.1",
				Host:  "example.com",
				Proto: "http",
			}},
			want: "by=192.168.1.1; for=203.0.113.1; host=example.com; proto=http",
		},
		{
			name: "2 for",
			forwarded: forwarded.Forwarded{
				{For: "203.0.113.1"},
				{For: "203.0.113.2"},
			},
			want: "for=203.0.113.1, for=203.0.113.2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.forwarded.StringSpace(); got != tt.want {
				t.Errorf("Forwarded.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestForwarded_Format(t *testing.T) {
	tests := []struct {
		format    string
		forwarded forwarded.Forwarded
		want      string
	}{
		{
			format: "%s",
			forwarded: forwarded.Forwarded{{
				By:  "192.168.1.1",
				For: "203.0.113.1",
			}},
			want: "by=192.168.1.1;for=203.0.113.1",
		},
		{
			format: "% s",
			forwarded: forwarded.Forwarded{{
				By:  "192.168.1.1",
				For: "203.0.113.1",
			}},
			want: "by=192.168.1.1; for=203.0.113.1",
		},
		{
			format: "% q",
			forwarded: forwarded.Forwarded{{
				For:  "203.0.113.1",
				Host: "the \"test\" server",
			}},
			want: "\"for=203.0.113.1; host=\\\"the \\\\\\\"test\\\\\\\" server\\\"\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			if got := fmt.Sprintf(tt.format, &tt.forwarded); got != tt.want {
				t.Errorf("Forwarded.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
