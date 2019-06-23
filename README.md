# Forwarded

Forwarded HTTP header ([RFC7239](https://tools.ietf.org/html/rfc7239)).

```sh
go get github.com/utisam/go-forwarded
```

```go
f := forwarded.Parse("for=192.0.2.43, for=\"[2001:db8:cafe::17]\", for=unknown")
fmt.Printf("%s\n", f)
fmt.Printf("% s\n", f) // With spaces
// Output:
// for=192.0.2.43,for=[2001:db8:cafe::17],for=unknown
// for=192.0.2.43, for=[2001:db8:cafe::17], for=unknown
```
