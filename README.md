# go-compose

Generic function chain composition in golang.

Turns `hf1(hf2(hf3(h)))` -> `compose.New(hf1, hf2, hf3).Compose(h)`

## Example Usage

### http.Handler Middleware Chaining

Turns:
```
auth(recoverPanic(logRequest(commonHeaders(mux))))
```

Into:
```
compose.New(auth, recoverPanic, logRequest, commonHeaders).Compose(mux)
```

### Math

Turns:
```
math.Floor(math.Sqrt(math.Abs(-1234)))
```

Into:
```
compose.New(math.Floor, math.Sqrt, math.Abs).Compose(-1234)
```