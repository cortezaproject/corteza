package helpers

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/textproto"
	"testing"

	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/require"
)

func InitFileUpload(t *testing.T, apiTest *apitest.APITest, endpoint string, form map[string]string, file []byte, name, mimetype string) *apitest.Response {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for k, v := range form {
		require.NoError(t, writer.WriteField(k, v))
	}

	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			"upload", name))
	header.Set("Content-Type", mimetype)
	part, err := writer.CreatePart(header)
	require.NoError(t, err)

	_, err = part.Write(file)
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	return apiTest.
		Post(endpoint).
		Body(body.String()).
		Header("Content-Type", writer.FormDataContentType()).
		Header("Accept", "application/json").
		Expect(t)
}
