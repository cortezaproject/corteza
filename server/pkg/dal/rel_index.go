package dal

import (
	"fmt"

	"github.com/spf13/cast"
)

type (
	// relIndex is a generic struct for indexing data which join/link can use
	//
	// The current index implementation utilizes a series of hashmaps based on what
	// type of values we're indexing on.
	// As we currently only support single value predicate, the hashmaps proved
	// to be a bit more efficient then b-trees (other consideration) when testing
	// on larger datasets.
	//
	// @todo when we support multiple join predicates, this should probably change
	// into a b-tree as it might be a bit faster thennested hashmaps.
	//
	// @todo do some benchmarks in regards to using generics for key
	relIndex struct {
		track []string

		keyType Type

		ints    map[int64]*relIndexBuffer
		strings map[string]*relIndexBuffer
		ids     map[uint64]*relIndexBuffer
	}
)

// newRelIndex initializes a new relIndex with the specified tracked attributes
// @todo benchmark with generics as key
func newRelIndex(t Type, track ...string) (out *relIndex, err error) {
	out = &relIndex{
		track:   track,
		keyType: t,
	}

	return out, out.initHashmaps()
}

// Add adds a new row to the index under the specified key
func (ri *relIndex) Add(k any, r *Row) {
	switch ri.keyType.(type) {
	case TypeNumber, *TypeNumber:
		ri.addInt(cast.ToInt64(k), r)
		return

	case TypeText, *TypeText:
		ri.addString(cast.ToString(k), r)
		return

	case TypeID, *TypeID,
		TypeRef, *TypeRef:
		ri.addID(cast.ToUint64(k), r)
		return
	}

	// @note this is validated when initializing
	panic(fmt.Sprintf("cannot use type %s as index key", ri.keyType.Type()))
}

func (ri *relIndex) Get(k any) (out *relIndexBuffer, ok bool) {
	switch ri.keyType.(type) {
	case TypeNumber, *TypeNumber:
		return ri.getInt(cast.ToInt64(k))

	case TypeText, *TypeText:
		return ri.getString(cast.ToString(k))

	case TypeID, *TypeID,
		TypeRef, *TypeRef:
		return ri.getID(cast.ToUint64(k))
	}

	// @note this is validated when initializing
	panic(fmt.Sprintf("cannot use type %s as index key", ri.keyType.Type()))
}

// Clear clears out the index omitting the need to reinitialize
func (ri *relIndex) Clear() (err error) {
	return ri.initHashmaps()
}

func (ri *relIndex) initHashmaps() (err error) {
	// @note initHashmaps initializes only the ones we need to save up on space
	switch ri.keyType.(type) {
	case TypeNumber, *TypeNumber:
		ri.ints = make(map[int64]*relIndexBuffer, 512)
		return

	case TypeText, *TypeText:
		ri.strings = make(map[string]*relIndexBuffer, 512)
		return

	case TypeID, *TypeID,
		TypeRef, *TypeRef:
		ri.ids = make(map[uint64]*relIndexBuffer, 512)
		return
	}
	return fmt.Errorf("cannot use type %s as index key", ri.keyType.Type())
}

func (ri *relIndex) addInt(k int64, r *Row) {
	c, ok := ri.getInt(k)
	if !ok {
		c = newRelIndexBuffer(ri.track...)
		ri.ints[k] = c
	}
	c.add(r)
}

func (ri *relIndex) getInt(k int64) (out *relIndexBuffer, ok bool) {
	out, ok = ri.ints[k]
	return
}

func (ri *relIndex) addString(k string, r *Row) {
	c, ok := ri.getString(k)
	if !ok {
		c = newRelIndexBuffer(ri.track...)
		ri.strings[k] = c
	}
	c.add(r)
}

func (ri *relIndex) getString(k string) (out *relIndexBuffer, ok bool) {
	out, ok = ri.strings[k]
	return
}

func (ri *relIndex) addID(k uint64, r *Row) {
	c, ok := ri.getID(k)
	if !ok {
		c = newRelIndexBuffer(ri.track...)
		ri.ids[k] = c
	}
	c.add(r)
}

func (ri *relIndex) getID(k uint64) (out *relIndexBuffer, ok bool) {
	out, ok = ri.ids[k]
	return
}
