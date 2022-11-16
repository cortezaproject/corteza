package workflows

import (
	"context"
	"io"
	"io/ioutil"
	"testing"

	"github.com/cortezaproject/corteza/server/automation/types"
	cmpTypes "github.com/cortezaproject/corteza/server/compose/types"
	"github.com/stretchr/testify/require"
)

func Test_attachment_management(t *testing.T) {
	var (
		ctx = bypassRBAC(context.Background())
		req = require.New(t)
	)

	req.NoError(defStore.TruncateAttachments(ctx))

	loadNewScenario(ctx, t)

	var (
		aux = struct {
			Base64blackGif    string
			LoadedAttBlackGif *cmpTypes.Attachment
			StoredAttBlackGif *cmpTypes.Attachment
			LoadedContent     io.ReadSeeker
		}{}
	)

	vars, _ := mustExecWorkflow(ctx, t, "attachments", types.WorkflowExecParams{})
	req.NoError(vars.Decode(&aux))

	req.Equal(int(aux.LoadedAttBlackGif.Meta.Original.Size), len(aux.Base64blackGif))
	req.Equal(int(aux.StoredAttBlackGif.Meta.Original.Size), len(aux.Base64blackGif))

	b, err := ioutil.ReadAll(aux.LoadedContent)
	req.NoError(err)
	req.Equal(aux.Base64blackGif, string(b))

}
