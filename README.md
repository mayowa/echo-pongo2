# echo-pongo2

Use pongo2 templates in the [echo](https://github.com/labstack/echo) web framework

## Usage

echo-pongo2 implements the echo [Renderer](http://godoc.org/github.com/labstack/echo#Renderer) interface

```go
e := echo.New()
r, err := echo-pongo2.NewRenderer("./template")
e.SetRenderer(r)
```

somewhere in a handler

```
func Hello(c *echo.Context) error {
    data := map[string]string{"World":"mayowa"}
    return c.Render(http.StatusOK, "hello.html", data)
}
```

template: ./template/hello.html

```
Hello {{World}}
```

## Options

```go
opts := echo-pongo2.Options {
    Debug: False,
    Source: echo-pongo2.FromFile
}
r, err := echo-pongo2.NewRenderer("./template", opts)
```

- Debug : When Debug is enabled Pongo2 wont cache parsed files
- Source [FromFile|FromCache]: Determines if subsequent request for a previously rendered template
  is retrieved from cache or if the template is re-read from file Cache. Defaults to FromCache
