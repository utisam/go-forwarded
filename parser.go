package forwarded

import (
	"errors"
	"net/http"
	"net/textproto"
	"strconv"
	"strings"
)

// ParseOption is a option to configure parser.
type ParseOption func(*parser)

// InitialCapacity is a option of perser to configure initial capacity of Forwarded slice.
func InitialCapacity(n int) ParseOption {
	return func(p *parser) {
		p.initCap = n
	}
}

// Parse forwarded field
func Parse(forwarded string, opts ...ParseOption) (Forwarded, error) {
	p := &parser{forwarded: forwarded}
	for _, opt := range opts {
		opt(p)
	}
	return p.parse()
}

// ErrNoHeaderFound is an error when no forwarded header found.
var ErrNoHeaderFound = errors.New("no forwarded header found")

// ParseHeader parse forwarded header fields
func ParseHeader(header http.Header, opts ...ParseOption) (Forwarded, error) {
	fields, ok := header[textproto.CanonicalMIMEHeaderKey("forwarded")]
	if !ok {
		return nil, ErrNoHeaderFound
	}

	res := Forwarded{}
	for _, field := range fields {
		elms, err := Parse(field)
		if err != nil {
			return nil, err
		}
		res = append(res, elms...)
	}
	return res, nil
}

type parser struct {
	forwarded string
	curent    int
	initCap   int
}

type paramToken string

const (
	tokenBy    = paramToken("by")
	tokenFor   = paramToken("for")
	tokenHost  = paramToken("host")
	tokenProto = paramToken("proto")
)

type pair struct {
	token paramToken
	value string
}

func (p *parser) parse() (Forwarded, error) {
	if len(p.forwarded) == 0 {
		return Forwarded{}, nil
	}

	res := make(Forwarded, 0, p.initCap)
	for {
		elem, err := p.parseElement()
		if err != nil {
			return res, err
		}

		res = append(res, elem)

		p.skipWitespace()
		if len(p.forwarded) <= p.curent || p.forwarded[p.curent] != ',' {
			break
		}
		p.curent++
	}

	return res, nil
}

func (p *parser) parseElement() (*Element, error) {
	res := Element{}

	for {
		pos := p.curent
		pair, err := p.parsePair()
		if err != nil {
			return nil, err
		}

		switch pair.token {
		case tokenBy:
			res.By = pair.value
		case tokenFor:
			res.For = pair.value
		case tokenHost:
			res.Host = pair.value
		case tokenProto:
			res.Proto = pair.value
		default:
			return nil, newParserError(pos, "unknown token: "+string(pair.token))
		}

		p.skipWitespace()
		if len(p.forwarded) <= p.curent || p.forwarded[p.curent] != ';' {
			break
		}
		p.curent++
	}

	return &res, nil
}

func (p *parser) parsePair() (*pair, error) {
	res := &pair{
		token: paramToken(strings.ToLower(p.readToken(func(delim byte) bool {
			return delim != '='
		}))),
	}

	p.skipWitespace()
	p.curent++ // '='

	p.skipWitespace()
	if len(p.forwarded) <= p.curent {
		return res, newParserError(p.curent, "unexpected EOS")
	}
	if p.forwarded[p.curent] == '"' {
		res.value = p.readQuotedString()
	} else {
		res.value = p.readToken(func(delim byte) bool {
			return delim != ';' && delim != ','
		})
	}

	return res, nil
}

func (p *parser) readQuotedString() string {
	p.curent++ // '"'

	b := strings.Builder{}
	// Maximum: IPv4-mapped IPv6 address and brackets
	b.Grow((4*6 + 5) + 1 + (3*4 + 3) + 2)

	escaped := false
	for p.curent < len(p.forwarded) && (escaped || p.forwarded[p.curent] != '"') {
		escaped = false
		if p.forwarded[p.curent] == '\\' {
			escaped = true
		} else {
			b.WriteByte(p.forwarded[p.curent])
		}
		p.curent++
	}

	p.curent++ // '"'

	return b.String()
}

func (p *parser) readToken(f func(delim byte) bool) string {
	p.skipWitespace()
	begin := p.curent
	for p.curent < len(p.forwarded) && !p.isWhitespace() && f(p.forwarded[p.curent]) {
		p.curent++
	}
	return p.forwarded[begin:p.curent]
}

func (p *parser) skipWitespace() {
	for p.curent < len(p.forwarded) && p.isWhitespace() {
		p.curent++
	}
}

func (p *parser) isWhitespace() bool {
	return p.forwarded[p.curent] == ' '
}

type parserError struct {
	position int
	message  string
}

func newParserError(pos int, msg string) error {
	return &parserError{pos, msg}
}

func (err *parserError) Error() string {
	return strconv.Itoa(err.position+1) + ": " + err.message
}
