package messaging

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"github.com/stretchr/testify/assert"

	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
)

func (h helper) apiChAttach(ch *types.Channel, file []byte) *apitest.Response {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("upload", "test.txt")
	h.a.NoError(err)

	_, err = part.Write(file)
	h.a.NoError(err)
	h.a.NoError(writer.Close())

	return h.apiInit().
		Post(fmt.Sprintf("/channels/%v/attach", ch.ID)).
		Body(body.String()).
		ContentType(writer.FormDataContentType()).
		Expect(h.t).
		Status(http.StatusOK)
}

// Non members should not be able to attach files to non-public channels
func TestChannelAttachNotMember(t *testing.T) {
	h := newHelper(t)

	ch := h.repoMakePrivateCh()

	h.apiChAttach(ch, []byte("NOPE")).
		Assert(helpers.AssertError("not allowed to attach files this channel")).
		End()
}

func TestChannelAttach(t *testing.T) {
	h := newHelper(t)

	uploadFileContent := "hello corteza, time here is " + time.Now().String()

	ch := h.repoMakePublicCh()
	h.repoMakeMember(ch, h.cUser)

	out := h.apiChAttach(ch, []byte(uploadFileContent)).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present(`$.response.url`)).
		Assert(jsonpath.Present(`$.response.previewUrl`)).
		End()

	// Now, fetch uploaded file and check the content
	var rval = struct {
		Response struct {
			Url        string
			PreviewUrl string
		}
	}{}
	out.JSON(&rval)
	attUrl, err := url.Parse(rval.Response.Url)
	assert.NoError(t, err)

	h.apiInit().
		Get(attUrl.Path).
		Query("sign", attUrl.Query().Get("sign")).
		Query("userID", attUrl.Query().Get("userID")).
		Expect(t).
		Status(http.StatusOK).
		Body(uploadFileContent).
		End()
}

func TestChannelAttachAndDelete(t *testing.T) {
	h := newHelper(t)

	ch := h.repoMakePublicCh()
	h.repoMakeMember(ch, h.cUser)

	h.apiChAttach(ch, []byte("dummy")).
		Assert(helpers.AssertNoErrors).
		End()

	var rval = struct {
		Response []struct {
			MessageID string
		}
	}{}

	h.apiInit().
		Get("/search/messages").
		Query("channelID", fmt.Sprintf("%d", ch.ID)).
		Expect(t).
		Assert(helpers.AssertNoErrors).
		End().
		JSON(&rval)

	h.a.Len(rval.Response, 1)

	h.apiInit().
		Delete(fmt.Sprintf("/channels/%d/messages/%s", ch.ID, rval.Response[0].MessageID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

}
