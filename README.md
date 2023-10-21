# Go Pool Package

The Go Pool package provides a pool implementation for managing preallocated memory in a concurrent environment. It is designed to efficiently manage a global rate limit.

```go
type Memory struct {
	// variables that you will use in handler.
}

p := pool.New[Value](2, time.Second) // only two request per second.

http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    err := p.With(r.Context(), func(v *Memory) {
        // v is preallocated.
        // [...]
    })
    if err != nil {
        // [handle error...]
    }
})

http.ListenAndServe(":8080", nil)
```
