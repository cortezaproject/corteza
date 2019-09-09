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

func chAttach(h helper, ch *types.Channel, file []byte) *apitest.Response {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("upload", "test.txt")
	h.a.NoError(err)

	_, err = part.Write(file)
	h.a.NoError(err)
	h.a.NoError(writer.Close())

	return h.testAPI().
		Debug().
		Post(fmt.Sprintf("/channels/%v/attach", ch.ID)).
		Body(body.String()).
		ContentType(writer.FormDataContentType()).
		Expect(h.t).
		Status(http.StatusOK)
}

// Non members should not be able to attach files to non-public channels
func TestChannelAttachNotMember(t *testing.T) {
	h := newHelper(t)

	ch := h.makePrivateCh()

	chAttach(h, ch, []byte("NOPE")).
		Assert(helpers.AssertError("messaging.service.NoPermissions")).
		End()
}

func TestChannelAttach(t *testing.T) {
	h := newHelper(t)

	uploadFileContent := "hello corteza, time here is " + time.Now().String()

	ch := h.makePublicCh()
	h.makeMember(ch, h.cUser)

	out := chAttach(h, ch, []byte(uploadFileContent)).
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

	h.testAPI().
		Get(attUrl.Path).
		Query("sign", attUrl.Query().Get("sign")).
		Query("userID", attUrl.Query().Get("userID")).
		Expect(t).
		Status(http.StatusOK).
		Body(uploadFileContent).
		End()
}
