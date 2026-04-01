package gmail

import (
	"context"

	"github.com/cortezaproject/corteza/server/store/adapters/api/drivers/google"
)

type gmailWrapper struct {
	*google.Wrapper
}

func newWrapper(baseURL string, connectionID uint64) *gmailWrapper {
	return &gmailWrapper{
		Wrapper: google.NewWrapper(baseURL, connectionID),
	}
}

func (w *gmailWrapper) Run(ctx context.Context, method string, path string, payload []byte, extraHeaders map[string][]string) (int, map[string][]string, []byte, error) {
	return w.Wrapper.Run(ctx, method, path, payload, extraHeaders)
}
