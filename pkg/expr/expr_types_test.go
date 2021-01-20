package expr

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"testing"
)

func TestKV_Set(t *testing.T) {
	var (
		req = require.New(t)

		vars = KV{value: map[string]string{
			"k1": "v1",
			"k2": "v2",
		}}
	)

	req.NoError(Assign(&vars, "k1", "v11"))
	req.Equal("v11", vars.value["k1"])
	req.Equal("v2", vars.value["k2"])

}

func TestKVV_Set(t *testing.T) {
	var (
		req = require.New(t)
		kvv KVV
	)

	req.NoError(Assign(&kvv, "foo", "bar"))
	req.Contains(kvv.value, "foo")
	req.Equal([]string{"bar"}, kvv.value["foo"])

	// Making sure http.Header is properly converted
	kvv = KVV{}
	req.NoError(kvv.Assign(http.Header{"foo": []string{"bar"}}))
	req.Contains(kvv.value, "foo")
	req.Equal([]string{"bar"}, kvv.value["foo"])

	// Making sure url.Values are properly converted
	kvv = KVV{}
	req.NoError(kvv.Assign(url.Values{"foo": []string{"bar"}}))
	req.Contains(kvv.value, "foo")
	req.Equal([]string{"bar"}, kvv.value["foo"])
}

func TestCastToArray(t *testing.T) {
	var (
		req = require.New(t)
	)

	arr, err := CastToArray([]interface{}{
		Must(NewString("abc")),
		Must(NewString("123")),
	})

	req.NoError(err)
	req.Len(arr, 2)
}

func TestArrayDecode(t *testing.T) {
	var (
		req = require.New(t)

		foo = struct {
			Typed   TypedValue
			Iface   interface{}
			Strings []string
			Values  []TypedValue
		}{}
	)

	arr, err := CastToArray([]interface{}{
		Must(NewString("abc")),
		Must(NewString("123")),
	})

	req.NoError(err)
	req.NoError(RVars{
		"strings": &Array{arr},
		"iface":   Must(NewString("typed")),
		"typed":   Must(NewString("typed")),
		"values":  &Array{arr},
	}.Vars().Decode(&foo))
	req.Len(foo.Strings, 2)
	req.Len(foo.Values, 2)
}
