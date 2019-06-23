package forwarded

import (
	"errors"
	"net/http"
)

//go:generate ./xfields.sh

var ErrInvalidLength = errors.New("invalid length")

func FromX(fields ...XField) (Forwarded, error) {
	if len(fields) == 0 {
		return Forwarded{}, nil
	}

	size := fields[0].len()
	for _, f := range fields[1:] {
		if f.len() != size {
			return nil, ErrInvalidLength
		}
	}

	res := make(Forwarded, size)
	for i := range res {
		elm := &Element{}
		for _, f := range fields {
			f.apply(elm, i)
		}
		res[i] = elm
	}

	return res, nil
}

func AlignX(fields ...XField) Forwarded {
	if len(fields) == 0 {
		return Forwarded{}
	}

	size := fields[0].len()
	for _, f := range fields[1:] {
		n := f.len()
		if size < n {
			size = n
		}
	}

	res := make(Forwarded, size)
	for i := range res {
		res[i] = &Element{}
	}

	for _, f := range fields {
		n := f.len()
		for i := 0; i < n; i++ {
			f.apply(res[i], i)
		}
	}

	return res
}

func AlignAllX(header http.Header) Forwarded {
	return AlignX(
		By(header.Get("X-Forwarded-By")),
		For(header.Get("X-Forwarded-For")),
		Host(header.Get("X-Forwarded-Host")),
		Proto(header.Get("X-Forwarded-Proto")),
		RealHost(header.Get("X-Real-Host")),
		RealIP(header.Get("X-Real-IP")),
	)
}
