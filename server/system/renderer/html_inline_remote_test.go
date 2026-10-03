package renderer

import (
	"bytes"
	"context"
	"errors"
	"testing"

	httpClient "github.com/crusttech/human/server/pkg/http"
)

// Templates are user content; inlineRemote must not reach addresses the
// restricted HTTP client refuses (cloud metadata, link-local)
func TestInlineRemoteRestricted(t *testing.T) {
	d := &genericHTMLDriver{}

	pl := &driverPayload{
		Template:  bytes.NewBufferString(`<img src="{{ inlineRemote "http://169.254.169.254/latest/meta-data/" }}">`),
		Variables: map[string]interface{}{},
	}

	_, err := d.Render(context.Background(), pl)
	if err == nil {
		t.Fatal("expected the metadata address to be refused")
	}

	if !errors.Is(err, httpClient.ErrRestrictedNetwork) && !bytes.Contains([]byte(err.Error()), []byte("restricted network")) {
		t.Fatalf("expected a restricted network error, got: %v", err)
	}
}
