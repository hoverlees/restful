# Introduction
This is A simple RESTful server written in GO. 

# Example
- Import this module
```go
import (
    "github.com/hoverlees/restful"
)
```

- Create server and write the handlers.
```go
server := restful.NewServer(":801")
server.AddRestfulHandler("/servers", http.MethodPost, "/{id}/info", func(w http.ResponseWriter, r *http.Request, uriParams map[string]string) {
    w.Write([]byte(fmt.Sprintf("server id is %s", uriParams["id"])))
})
server.AddRestfulHandler("/servers", http.MethodGet, "/{id}/info/{field}", func(w http.ResponseWriter, r *http.Request, uriParams map[string]string) {
    w.Write([]byte(fmt.Sprintf("get %s for %s", uriParams["field"], uriParams["id"])))
})
server.Start()
```
