package echopongo2

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/flosch/pongo2/v4"
)

func TestRenderer(t *testing.T) {
	_, err := NewRenderer("/tmp")
	if err != nil {
		t.Error(err)
	}
}

func TestRender(t *testing.T) {
	baseDir := "/tmp"
	tpl, err := NewRenderer(baseDir)
	if err != nil {
		t.Error(err)
	}

	tplNme, err := makeTemplate(baseDir)
	if err != nil {
		t.Error(err)
	}

	buff := bytes.Buffer{}
	err = tpl.Render(&buff, tplNme, map[string]string{"World": "mayowa"}, nil)
	if err != nil {
		t.Error(err)
	}

	if buff.String() != "Hello mayowa!" {
		t.Errorf("Template not properly rendered: got ==> %s", buff.String())
	}

}

func TestRenderWithDebug(t *testing.T) {
	baseDir := "/tmp"
	tpl, err := NewRenderer(baseDir, Options{Debug: true})
	if err != nil {
		t.Error(err)
		return
	}

	tplNme, err := makeTemplate(baseDir)
	if err != nil {
		t.Error(err)
		return
	}

	buff := bytes.Buffer{}
	err = tpl.Render(&buff, tplNme, map[string]string{"World": "mayowa"}, nil)
	if err != nil {
		t.Error(err)
		return
	}

	if buff.String() != "Hello mayowa!" {
		t.Errorf("Template not properly rendered: got ==> %s", buff.String())
		return
	}

	err = modifyTemplate(baseDir, tplNme, "jumping {{World}}!")
	if err != nil {
		t.Error(err)
	}

	buff = bytes.Buffer{}
	err = tpl.Render(&buff, tplNme, map[string]string{"World": "mayowa"}, nil)
	if err != nil {
		t.Error(err)
		return
	}

	if buff.String() != "jumping mayowa!" {
		t.Errorf("Template not properly rendered: got ==> %s", buff.String())
		return
	}

}

func TestRenderWithSource(t *testing.T) {
	baseDir := "/tmp"
	tpl, err := NewRenderer(baseDir, Options{Source: FromFile})
	if err != nil {
		t.Error(err)
		return
	}

	tplNme, err := makeTemplate(baseDir)
	if err != nil {
		t.Error(err)
		return
	}

	buff := bytes.Buffer{}
	err = tpl.Render(&buff, tplNme, map[string]string{"World": "mayowa"}, nil)
	if err != nil {
		t.Error(err)
		return
	}

	if buff.String() != "Hello mayowa!" {
		t.Errorf("Template not properly rendered: got ==> %s", buff.String())
		return
	}

	err = modifyTemplate(baseDir, tplNme, "jumping {{World}}!")
	if err != nil {
		t.Error(err)
	}

	buff = bytes.Buffer{}
	err = tpl.Render(&buff, tplNme, map[string]string{"World": "mayowa"}, nil)
	if err != nil {
		t.Error(err)
		return
	}

	if buff.String() != "jumping mayowa!" {
		t.Errorf("Template not properly rendered: got ==> %s", buff.String())
		return
	}
}

func TestRenderWithoutSource(t *testing.T) {
	baseDir := "/tmp"
	tpl, err := NewRenderer(baseDir)
	if err != nil {
		t.Error(err)
		return
	}

	tplNme, err := makeTemplate(baseDir)
	if err != nil {
		t.Error(err)
		return
	}

	buff := bytes.Buffer{}
	err = tpl.Render(&buff, tplNme, map[string]string{"World": "mayowa"}, nil)
	if err != nil {
		t.Error(err)
		return
	}

	if buff.String() != "Hello mayowa!" {
		t.Errorf("Template not properly rendered: got ==> %s", buff.String())
		return
	}

	err = modifyTemplate(baseDir, tplNme, "jumping {{World}}!")
	if err != nil {
		t.Error(err)
	}

	buff = bytes.Buffer{}
	err = tpl.Render(&buff, tplNme, map[string]string{"World": "mayowa"}, nil)
	if err != nil {
		t.Error(err)
		return
	}

	if buff.String() != "Hello mayowa!" {
		t.Errorf("Template not properly rendered: got ==> %s", buff.String())
		return
	}
}

// make a on disk template and return its name
func makeTemplate(baseDir string) (string, error) {
	tplStr := `Hello {{World}}!`
	fNme := filepath.Join(baseDir, "test1.html")
	fHdl, err := os.Create(fNme)
	if err != nil {
		return "", err
	}
	defer fHdl.Close()

	_, err = fHdl.WriteString(tplStr)
	if err != nil {
		return "", err
	}

	return "test1.html", nil
}

func modifyTemplate(baseDir, name, content string) error {

	fNme := filepath.Join(baseDir, name)
	err := ioutil.WriteFile(fNme, []byte(content), 0x777)
	if err != nil {
		return err
	}

	return nil
}

func TestToPongoCtx(t *testing.T) {
	// test pongo2.Context
	v := pongo2.Context{"a": 1, "b": 2, "c": 3}
	retv, err := toPongoCtx(v)
	if err != nil {
		t.Error(err)
	}
	if retv["a"] != 1 || retv["b"] != 2 || retv["c"] != 3 {
		t.Errorf("Input data was mangled: is %v should be %v", retv, v)
	}

	// test simple struct
	type TStruct struct {
		A, B, C int
	}
	s := TStruct{A: 1, B: 2, C: 3}
	retv, err = toPongoCtx(s)
	if err != nil {
		t.Error(err)
	}
	if retv["A"] != 1 || retv["B"] != 2 || retv["C"] != 3 {
		t.Errorf("[Simple Struct]Input data was mangled: is %v should be %v", retv, s)
	}

	// test nested struct
	type TNested struct {
		A int
		B TStruct
	}
	n := TNested{A: -1, B: TStruct{A: 1, B: 2, C: 3}}
	retv, err = toPongoCtx(n)
	if err != nil {
		t.Error(err)
	}
	nb := retv["B"].(TStruct)
	if retv["A"] != -1 || nb.A != 1 || nb.B != 2 || nb.C != 3 {
		t.Errorf("[Nested Struct] Input data was mangled: is %v should be %v", retv, n)
	}

	// test map[string]int
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	retv, err = toPongoCtx(m)
	if err != nil {
		t.Error(err)
	}
	if retv["a"] != 1 || retv["b"] != 2 || retv["c"] != 3 {
		t.Errorf("[Map-String-Int] Input data was mangled: is %v should be %v", retv, m)
	}

	// test map[string]string
	m2 := map[string]string{"a": "1", "b": "2", "c": "3"}
	retv, err = toPongoCtx(m2)
	if err != nil {
		t.Error(err)
	}
	if retv["a"] != "1" || retv["b"] != "2" || retv["c"] != "3" {
		t.Errorf("[Map-String-Int] Input data was mangled: is %v should be %v", retv, m2)
	}

	// test map[int]string
	m3 := map[int]string{1: "1", 2: "2", 3: "3"}
	retv, err = toPongoCtx(m3)
	if err == nil {
		t.Error(fmt.Errorf("only map[string] is supported. An error should have been returned: %v", retv))
	}

}
