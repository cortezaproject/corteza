package renderer

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"
)

func TestGenericHTMLHeaderFooter(t *testing.T) {
	d := &genericHTMLDriver{}

	pl := &driverPayload{
		Template:  bytes.NewBufferString(`<main>{{template "pp" .}}</main>`),
		Header:    bytes.NewBufferString(`<header>{{.name}} {{template "pp" .}}</header>`),
		Footer:    bytes.NewBufferString(`<footer>{{.name}}</footer>`),
		Variables: map[string]interface{}{"name": "Human"},
		Partials: map[string]io.Reader{
			"pp": bytes.NewBufferString(`{{define "pp"}}partial:{{.name}}{{end}}`),
		},
	}

	out, err := d.Render(context.Background(), pl)
	if err != nil {
		t.Fatal(err)
	}

	bb, _ := io.ReadAll(out)
	got := string(bb)
	want := `<header>Human partial:Human</header><main>partial:Human</main><footer>Human</footer>`
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}

	// no header/footer keeps the single-pass path working
	pl2 := &driverPayload{
		Template:  bytes.NewBufferString(`<main>ok</main>`),
		Variables: map[string]interface{}{},
	}
	out2, err := d.Render(context.Background(), pl2)
	if err != nil {
		t.Fatal(err)
	}
	bb2, _ := io.ReadAll(out2)
	if !strings.Contains(string(bb2), "<main>ok</main>") {
		t.Fatalf("plain render broken: %q", string(bb2))
	}
}

func TestGenericHTMLSafeHTML(t *testing.T) {
	d := &genericHTMLDriver{}
	rte := `<h1>Title 1</h1>`

	pl := &driverPayload{
		Template:  bytes.NewBufferString(`<div>{{.body}}</div><div>{{.body | safeHTML}}</div>`),
		Variables: map[string]interface{}{"body": rte},
	}

	out, err := d.Render(context.Background(), pl)
	if err != nil {
		t.Fatal(err)
	}

	bb, _ := io.ReadAll(out)
	got := string(bb)
	want := `<div>&lt;h1&gt;Title 1&lt;/h1&gt;</div><div><h1>Title 1</h1></div>`
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}

	// the case the partial workaround could not cover: rich text in a list
	pl2 := &driverPayload{
		Template:  bytes.NewBufferString(`{{range .items}}<td>{{.desc | safeHTML}}</td>{{end}}`),
		Variables: map[string]interface{}{"items": []map[string]any{{"desc": rte}, {"desc": rte}}},
	}

	out2, err := d.Render(context.Background(), pl2)
	if err != nil {
		t.Fatal(err)
	}

	bb2, _ := io.ReadAll(out2)
	if got, want := string(bb2), `<td><h1>Title 1</h1></td><td><h1>Title 1</h1></td>`; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}

	// an unset or empty value must not break the render of the whole document
	pl3 := &driverPayload{
		Template:  bytes.NewBufferString(`<td>{{.unset | safeHTML}}</td><td>{{.blank | safeHTML}}</td><td>{{.count | safeHTML}}</td>`),
		Variables: map[string]interface{}{"blank": nil, "count": 3},
	}

	out3, err := d.Render(context.Background(), pl3)
	if err != nil {
		t.Fatal(err)
	}

	bb3, _ := io.ReadAll(out3)
	if got, want := string(bb3), `<td></td><td></td><td>3</td>`; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
