package echopongo2

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"github.com/labstack/echo"
	"github.com/flosch/pongo2"
)

// Renderer manages a pongo2 TemplateSet
type Renderer struct {
	BaseDir string
	TplSet  *pongo2.TemplateSet
}

// NewRenderer creates a new instance of Renderer
func NewRenderer(baseDir string) (*Renderer, error) {
	// check if baseDir exists
	fInfo, err := os.Lstat(baseDir)
	if err != nil {
		return nil, err
	}
	if fInfo.IsDir() == false {
		return nil, fmt.Errorf("%s is not a directory", baseDir)
	}

	rdr := Renderer{}
	loader, err := pongo2.NewLocalFileSystemLoader(baseDir)
	if err != nil {
		return nil, err
	}

	rdr.TplSet = pongo2.NewSet("TplSet-"+filepath.Base(baseDir), loader)

	return &rdr, nil
}

// Render implements echo.Render interface
func (r *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	// get template, compile it anf store it in cache
	tpl, err := r.TplSet.FromCache(name)
	if err != nil {
		return err
	}
	// convert supplied data to pongo2.Context
	val, err := toPongoCtx(data)
	if err != nil {
		return err
	}
	// generate render the template
	err = tpl.ExecuteWriter(val, w)
	return err
}

// toPongoCtx converts a pongo2.Context, struct, map[string] to
// pongo2.Context
func toPongoCtx(data interface{}) (pongo2.Context, error) {
	m := pongo2.Context{}

	v := reflect.ValueOf(data)
	if v.Type().String() == "pongo2.Context" {
		return data.(pongo2.Context), nil
	} else if v.Kind().String() == "struct" {
		for i := 0; i < v.NumField(); i++ {
			m[v.Type().Field(i).Name] = v.Field(i).Interface()
		}
	} else if strings.HasPrefix(v.Type().String(), "map[string]") {
		for _, k := range v.MapKeys() {
			m[k.String()] = v.MapIndex(k).Interface()
		}
	} else {
		return nil, fmt.Errorf("cant convert %T to pongo2.Context", data)
	}

	return m, nil
}
