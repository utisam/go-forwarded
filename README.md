# Forwarded

[![GoDoc](https://godoc.org/github.com/utisam/go-forwarded?status.svg)](https://godoc.org/github.com/utisam/go-forwarded)

Forwarded HTTP header ([RFC7239](https://tools.ietf.org/html/rfc7239)).

```sh
go get github.com/utisam/go-forwarded
```

```go
import "github.com/utisam/go-forwarded"

f := forwarded.Parse("for=192.0.2.43, for=\"[2001:db8:cafe::17]\", for=unknown")
fmt.Printf("%s\n", f)
fmt.Printf("% s\n", f) // With spaces
// Output:
// for=192.0.2.43,for=[2001:db8:cafe::17],for=unknown
// for=192.0.2.43, for=[2001:db8:cafe::17], for=unknown
```
