package yaml

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestUser_UnmarshalYAML(t *testing.T) {
	var (
		parseString = func(src string) (*user, error) {
			w := &user{}
			return w, yaml.Unmarshal([]byte(src), w)
		}
	)

	t.Run("empty", func(t *testing.T) {
		req := require.New(t)

		w, err := parseString(``)
		req.NoError(err)
		req.NotNil(w)
		req.Nil(w.res)
	})

	t.Run("simple name", func(t *testing.T) {
		req := require.New(t)

		w, err := parseString(`{ name: Test }`)
		req.NoError(err)
		req.NotNil(w)
		req.NotNil(w.res)
		req.NotEmpty(w.res.Name)
	})

	t.Run("user 1", func(t *testing.T) {
		req := require.New(t)

		doc, err := parseDocument("user_1")
		req.NoError(err)
		req.NotNil(doc)
		req.Len(doc.users, 2)
		req.NotNil(doc.users[0])
		req.Equal("u1@example.tld", doc.users[0].res.Email)
		req.Equal("u1", doc.users[0].res.Handle)
		req.NotNil(doc.users[1])
		req.Equal("John Doe", doc.users[1].res.Name)
		req.Equal("jdoe@example.tld", doc.users[1].res.Email)
		req.Equal(true, doc.users[1].res.EmailConfirmed)
		req.Equal("u2", doc.users[1].res.Handle)
	})
}
