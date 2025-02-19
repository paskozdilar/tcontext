# tcontext

TContext is a type-safe wrapper around Go context.Context that provides a way
to store and retrieve values in a type-safe, transparent way.

## Usage

It is recommended to use a pointer to a struct that contains all relevant
fields that will be used in the context. This way, each field can be accessed
and modified in a simple way across all functions that use the context.

The `tcontext.Context` type implements `context.Context`, so it can be used
wherever a `context.Context` is expected.
It also stores the Data as a `context.Context` value, so it can be retrieved
from child context using `tcontext.FromContext()`.

E.g.

```go
tctx := tcontext.WithData(context.Background(), &RequestData{
    RequestID: "123",
    UserID: "456",
})

ctx, cancel := context.WithCancel(tctx)
defer cancel()

tctx, ok = tcontext.FromContext[*RequestData](ctx)
if ok {
    fmt.Println(tctx.Data().RequestID)
    fmt.Println(tctx.Data().UserID)
}
```

For a more elaborate example, see
[./examples/http/main.go](./examples/http/main.go).
