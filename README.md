# echo-pongo2
Use pongo2 templates in the [echo](https://github.com/labstack/echo) web framework

## Usage
echo-pongo2 implements the echo [Renderer](http://godoc.org/github.com/labstack/echo#Renderer) interface


```
e := echo.New()
r, err := echo-pongo2.NewRenderer("./template")
e.SetRenderer(r)
```

somewhere in a handler
```
func Hello(c *echo.Context) error {
    return c.Render(http.StatusOK, "hello.html", "World")
}
```

template: ./template/hello.html
```
Hello {{World}}
```
