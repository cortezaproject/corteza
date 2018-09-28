package store

import (
	"bytes"
	"io"
	"testing"
)

func TestStore(t *testing.T) {
	assert := func(ok bool, format string, params ...interface{}) {
		if !ok {
			t.Fatalf(format, params...)
		}
	}
	readerToString := func(r io.Reader) string {
		b := new(bytes.Buffer)
		b.ReadFrom(r)
		return b.String()
	}

	store, err := New("test")
	assert(err == nil, "Unexpected error when creating store: %+v", err)
	assert(store != nil, "Expected non-nil return for new store")
	assert(store.Namespace() == "test", "Unexpected store namespace: test != %s", store.Namespace())

	{
		fn := store.Original(123, "jpg")
		expected := "test/123.jpg"
		assert(fn == expected, "Unexpected filename returned: %s != %s", expected, fn)
	}

	{
		fn := store.Preview(123, "jpg")
		expected := "test/123_preview.jpg"
		assert(fn == expected, "Unexpected filename returned: %s != %s", expected, fn)
	}

	// write a file
	{
		buf := bytes.NewBuffer([]byte("This is a testing buffer"))
		err := store.Save("test/123.jpg", buf)
		assert(err == nil, "Error saving file, %+v", err)

		err = store.Save("test123/123.jpg", buf)
		assert(err != nil, "Expected error when saving file outside of namespace")
	}

	// read a file
	{
		buf, err := store.Open("test/123.jpg")
		assert(err == nil, "Unexpected error when reading file: %+v", err)
		s := readerToString(buf)
		assert(s == "This is a testing buffer", "Unexpected response when reading file: %s", s)

		_, err = store.Open("test/1234.jpg")
		assert(err != nil, "Expected error when opening non-existent file")
		_, err = store.Open("test123/123.jpg")
		assert(err != nil, "Expected error when opening file outside of namespace")
	}

	// delete a file
	{
		err := store.Remove("test/123.jpg")
		assert(err == nil, "Unexpected error when removing file: %+v", err)
		err = store.Remove("test/123.jpg")
		assert(err != nil, "Expected error when removing missing file")
		err = store.Remove("test123/123.jpg")
		assert(err != nil, "Expected error when deleting file outside of namespace")
	}
}
