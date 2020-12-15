package echopongo2

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"

	"github.com/flosch/pongo2/v4"
	"github.com/labstack/echo/v4"
)

// Renderer manages a pongo2 TemplateSet
type Renderer struct {
	BaseDir string
	TplSet  *pongo2.TemplateSet
	debug   bool
	source  RenderSource
	globals map[string]interface{}
}

// RenderSource source from which template will be rendered
type RenderSource int

const (
	// FromCache render template from cache
	FromCache RenderSource = iota
	// FromFile render template from file
	FromFile
)

// Options to modify the renders behavior
type Options struct {
	Debug             bool
	Source            RenderSource
	MixManifestFolder string
}

// NewRenderer creates a new instance of Renderer
func NewRenderer(baseDir string, opts ...Options) (*Renderer, error) {
	// check if baseDir exists
	fInfo, err := os.Lstat(baseDir)
	if err != nil {
		return nil, err
	}
	if fInfo.IsDir() == false {
		return nil, fmt.Errorf("%s is not a directory", baseDir)
	}

	rdr := Renderer{}

	if opts != nil {
		err := rdr.RegisterFilter("mix", MixManifest(opts[0].MixManifestFolder))
		if err != nil {
			return nil, err
		}
	}

	for _, i := range opts {
		rdr.debug = i.Debug
		rdr.source = i.Source
	}

	loader, err := pongo2.NewLocalFileSystemLoader(baseDir)
	if err != nil {
		return nil, err
	}

	rdr.TplSet = pongo2.NewSet("TplSet-"+filepath.Base(baseDir), loader)
	rdr.TplSet.Debug = rdr.debug

	// allocate globals map
	rdr.globals = make(map[string]interface{})

	return &rdr, nil
}

// Render implements echo.Render interface
func (r *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	var (
		tpl *pongo2.Template
		err error
	)

	// get template, compile it anf store it in cache
	if r.source == FromFile {
		tpl, err = r.TplSet.FromFile(name)
	} else {
		tpl, err = r.TplSet.FromCache(name)
	}
	if err != nil {
		return err
	}

	// convert supplied data to pongo2.Context
	val, err := toPongoCtx(data)
	if err != nil {
		return err
	}

	// add globals to new context
	for k, v := range r.globals {
		val[k] = v
	}

	// generate render the template
	err = tpl.ExecuteWriter(val, w)
	return err
}

// toPongoCtx converts a pongo2.Context, struct, map[string] to
// pongo2.Context
func toPongoCtx(data interface{}) (pongo2.Context, error) {
	m := pongo2.Context{}

	v := reflect.Indirect(reflect.ValueOf(data))
	if v.Type().String() == "pongo2.Context" {
		return data.(pongo2.Context), nil
	} else if v.Kind().String() == "struct" {
		for i := 0; i < v.NumField(); i++ {
			m[v.Type().Field(i).Name] = v.Field(i).Interface()
		}
	} else if v.Type().Kind() == reflect.Map && v.Type().Key().Kind() == reflect.String {
		for _, k := range v.MapKeys() {
			// fmt.Println("k:", k.String(), k)
			m[k.String()] = v.MapIndex(k).Interface()
		}
	} else {
		return nil, fmt.Errorf("cant convert %T to pongo2.Context", data)
	}

	return m, nil
}

// RegisterFilter ...
func (r *Renderer) RegisterFilter(name string, fn pongo2.FilterFunction) error {
	if pongo2.FilterExists(name) {
		return pongo2.ReplaceFilter(name, fn)
	}

	return pongo2.RegisterFilter(name, fn)
}

// SetGlobal ...
func (r *Renderer) SetGlobal(name string, val interface{}) {
	r.globals[name] = val
}
