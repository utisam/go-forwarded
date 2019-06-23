package forwarded

import (
	"fmt"
	"io"
	"strings"
)

// Forwarded is a slice of each step.
type Forwarded []*Element

// String returns a string value without extra spaces.
func (f Forwarded) String() string {
	ss := make([]string, len(f))
	for i := range f {
		ss[i] = f[i].String()
	}
	return strings.Join(ss, ",")
}

// StringSpace returns a string value with spaces.
func (f Forwarded) StringSpace() string {
	ss := make([]string, len(f))
	for i := range f {
		ss[i] = f[i].StringSpace()
	}
	return strings.Join(ss, ", ")
}

// Format is a method for fmt.Formatter.
func (f Forwarded) Format(s fmt.State, c rune) {
	var str string
	if s.Flag(' ') {
		str = f.StringSpace()
	} else {
		str = f.String()
	}

	switch c {
	case 'v', 's':
		io.WriteString(s, str)
	case 'q':
		fmt.Fprintf(s, "%q", str)
	}
}

// Element is a step of proxy.
type Element struct {
	By    string `json:"by,omitempty"`
	For   string `json:"for,omitempty"`
	Host  string `json:"host,omitempty"`
	Proto string `json:"proto,omitempty"`
}

// String returns a string value without extra spaces.
func (e *Element) String() string {
	return strings.Join(e.stringPairs(), ";")
}

// StringSpace returns a string value with spaces.
func (e *Element) StringSpace() string {
	return strings.Join(e.stringPairs(), "; ")
}

func (e *Element) stringPairs() []string {
	res := make([]string, 0, 4)
	res = appendPair(res, "by", e.By)
	res = appendPair(res, "for", e.For)
	res = appendPair(res, "host", e.Host)
	res = appendPair(res, "proto", e.Proto)
	return res
}

func appendPair(s []string, name, v string) []string {
	if len(v) != 0 {
		s = append(s, name+"="+stringValue(v))
	}
	return s
}

func stringValue(s string) string {
	quote := false
	escape := 0
	for _, r := range s {
		switch r {
		case '"':
			escape++
			fallthrough
		case ';', ',', ' ':
			quote = true
		}
	}

	if !quote {
		return s
	}

	b := strings.Builder{}
	b.Grow(len(s) + 2 /* quotes */ + escape)
	b.WriteRune('"')
	for _, r := range s {
		if r == '"' {
			b.WriteRune('\\')
		}
		b.WriteRune(r)
	}
	b.WriteRune('"')
	return b.String()
}
