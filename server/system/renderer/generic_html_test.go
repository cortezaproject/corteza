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
