# btree

[![GoDoc](https://godoc.org/github.com/tidwall/btree?status.svg)](https://godoc.org/github.com/tidwall/btree)

An efficient [B-tree](https://en.wikipedia.org/wiki/B-tree) implementation in Go.

## Features

- Support for [Generics](#generics) (Go 1.18+).
- `Map` and `Set` types for ordered key-value maps and sets,
- Fast bulk loading for pre-ordered data using the `Load()` method.
- `Copy()` method with copy-on-write support.
- Thread-safe operations.
- [Path hinting](PATH_HINT.md) optimization for operations with nearby keys.

## Using

To start using this package, install Go and run:

```sh
$ go get github.com/tidwall/btree
```

## B-tree types

This package includes the following types of B-trees:

- [`btree.Map`](#btreemap):
A fast B-tree for storing ordered key value pairs.
Go 1.18+ 
- [`btree.Set`](#btreeset):
Like `Map`, but only for storing keys.
Go 1.18+
- [`btree.Generic`](#btreegeneric):
A feature-rich B-tree for storing data using a custom comparator.
Go 1.18+
- [`btree.BTree`](#btreebtree):
Like `Generic` but uses the `interface{}` type for data. Backwards compatible.
Go 1.16+

### btree.Map

```go
// Basic
Set(key, value)    // insert or replace an item
Get(key, value)    // get an existing item
Delete(key)        // delete an item
Len()              // return the number of items in the map

// Iteration
Scan(iter)         // scan items in ascending order
Reverse(iter)      // scan items in descending order
Ascend(key, iter)  // scan items in ascending order that are >= to key
Descend(key, iter) // scan items in descending order that are <= to key.
Iter()             // returns a read-only iterator for for-loops.

// Array-like operations
GetAt(index)       // returns the item at index
DeleteAt(index)    // deletes the item at index

// Bulk-loading
Load(key, value)   // load presorted items into tree
```

#### Example

```go
package main

import (
	"fmt"
	"github.com/tidwall/btree"
)

func main() {
	// create a map
	var users btree.Map[string, string]

	// add some users
	users.Set("user:4", "Andrea")
	users.Set("user:6", "Andy")
	users.Set("user:2", "Andy")
	users.Set("user:1", "Jane")
	users.Set("user:5", "Janet")
	users.Set("user:3", "Steve")

	// Iterate over the maps and print each user
	users.Scan(func(key, value string) bool {
		fmt.Printf("%s %s\n", key, value)
		return true
	})
	fmt.Printf("\n")

	// Delete a couple
	users.Delete("user:5")
	users.Delete("user:1")

	// print the map again
	users.Scan(func(key, value string) bool {
		fmt.Printf("%s %s\n", key, value)
		return true
	})
	fmt.Printf("\n")

	// Output:
	// user:1 Jane
	// user:2 Andy
	// user:3 Steve
	// user:4 Andrea
	// user:5 Janet
	// user:6 Andy
	//
	// user:2 Andy
	// user:3 Steve
	// user:4 Andrea
	// user:6 Andy
}
```

### btree.Set

```go
// Basic
Insert(key)        // insert an item
Contains(key)      // test if item exists
Delete(key)        // delete an item
Len()              // return the number of items in the set

// Iteration
Scan(iter)         // scan items in ascending order
Reverse(iter)      // scan items in descending order
Ascend(key, iter)  // scan items in ascending order that are >= to key
Descend(key, iter) // scan items in descending order that are <= to key.
Iter()             // returns a read-only iterator for for-loops.

// Array-like operations
GetAt(index)       // returns the item at index
DeleteAt(index)    // deletes the item at index

// Bulk-loading
Load(key)          // load presorted item into tree
```

#### Example

```go
package main

import (
	"fmt"
	"github.com/tidwall/btree"
)

func main() {
	// create a set
	var names btree.Set[string]

	// add some names
	names.Insert("Jane")
	names.Insert("Andrea")
	names.Insert("Steve")
	names.Insert("Andy")
	names.Insert("Janet")
	names.Insert("Andy")

	// Iterate over the maps and print each user
	names.Scan(func(key string) bool {
		fmt.Printf("%s\n", key)
		return true
	})
	fmt.Printf("\n")

	// Delete a couple
	names.Delete("Steve")
	names.Delete("Andy")

	// print the map again
	names.Scan(func(key string) bool {
		fmt.Printf("%s\n", key)
		return true
	})
	fmt.Printf("\n")

	// Output:
	// Andrea
	// Andy
	// Jane
	// Janet
	// Steve
	//
	// Andrea
	// Jane
	// Janet
}
```

### btree.Generic

```go
// Basic
Set(item)               // insert or replace an item
Get(item)               // get an existing item
Delete(item)            // delete an item
Len()                   // return the number of items in the btree

// Iteration
Scan(iter)              // scan items in ascending order
Reverse(iter)           // scan items in descending order
Ascend(key, iter)       // scan items in ascending order that are >= to key
Descend(key, iter)      // scan items in descending order that are <= to key.
Iter()                  // returns a read-only iterator for for-loops.

// Array-like operations
GetAt(index)            // returns the item at index
DeleteAt(index)         // deletes the item at index

// Bulk-loading
Load(item)              // load presorted items into tree

// Path hinting
SetHint(item, *hint)    // insert or replace an existing item
GetHint(item, *hint)    // get an existing item
DeleteHint(item, *hint) // delete an item

// Copy-on-write
Copy()                  // copy the btree
```

#### Example

```go
package main

import (
	"fmt"

	"github.com/tidwall/btree"
)

type Item struct {
	Key, Val string
}

// byKeys is a comparison function that compares item keys and returns true
// when a is less than b.
func byKeys(a, b Item) bool {
	return a.Key < b.Key
}

// byVals is a comparison function that compares item values and returns true
// when a is less than b.
func byVals(a, b Item) bool {
	if a.Val < b.Val {
		return true
	}
	if a.Val > b.Val {
		return false
	}
	// Both vals are equal so we should fall though
	// and let the key comparison take over.
	return byKeys(a, b)
}

func main() {
	// Create a tree for keys and a tree for values.
	// The "keys" tree will be sorted on the Keys field.
	// The "values" tree will be sorted on the Values field.
	keys := btree.NewGeneric[Item](byKeys)
	vals := btree.NewGeneric[Item](byVals)

	// Create some items.
	users := []Item{
		Item{Key: "user:1", Val: "Jane"},
		Item{Key: "user:2", Val: "Andy"},
		Item{Key: "user:3", Val: "Steve"},
		Item{Key: "user:4", Val: "Andrea"},
		Item{Key: "user:5", Val: "Janet"},
		Item{Key: "user:6", Val: "Andy"},
	}

	// Insert each user into both trees
	for _, user := range users {
		keys.Set(user)
		vals.Set(user)
	}

	// Iterate over each user in the key tree
	keys.Scan(func(item Item) bool {
		fmt.Printf("%s %s\n", item.Key, item.Val)
		return true
	})
	fmt.Printf("\n")

	// Iterate over each user in the val tree
	vals.Scan(func(item Item) bool {
		fmt.Printf("%s %s\n", item.Key, item.Val)
		return true
	})

	// Output:
	// user:1 Jane
	// user:2 Andy
	// user:3 Steve
	// user:4 Andrea
	// user:5 Janet
	// user:6 Andy
	//
	// user:4 Andrea
	// user:2 Andy
	// user:6 Andy
	// user:1 Jane
	// user:5 Janet
	// user:3 Steve
}
```

### btree.BTree

```go
// Basic
Set(item)               // insert or replace an item
Get(item)               // get an existing item
Delete(item)            // delete an item
Len()                   // return the number of items in the btree

// Iteration
Scan(iter)              // scan items in ascending order
Reverse(iter)           // scan items in descending order
Ascend(key, iter)       // scan items in ascending order that are >= to key
Descend(key, iter)      // scan items in descending order that are <= to key.
Iter()                  // returns a read-only iterator for for-loops.

// Array-like operations
GetAt(index)            // returns the item at index
DeleteAt(index)         // deletes the item at index

// Bulk-loading
Load(item)              // load presorted items into tree

// Path hinting
SetHint(item, *hint)    // insert or replace an existing item
GetHint(item, *hint)    // get an existing item
DeleteHint(item, *hint) // delete an item

// Copy-on-write
Copy()                  // copy the btree
```

#### Example

```go
package main

import (
	"fmt"

	"github.com/tidwall/btree"
)

type Item struct {
	Key, Val string
}

// byKeys is a comparison function that compares item keys and returns true
// when a is less than b.
func byKeys(a, b interface{}) bool {
	i1, i2 := a.(*Item), b.(*Item)
	return i1.Key < i2.Key
}

// byVals is a comparison function that compares item values and returns true
// when a is less than b.
func byVals(a, b interface{}) bool {
	i1, i2 := a.(*Item), b.(*Item)
	if i1.Val < i2.Val {
		return true
	}
	if i1.Val > i2.Val {
		return false
	}
	// Both vals are equal so we should fall though
	// and let the key comparison take over.
	return byKeys(a, b)
}

func main() {
	// Create a tree for keys and a tree for values.
	// The "keys" tree will be sorted on the Keys field.
	// The "values" tree will be sorted on the Values field.
	keys := btree.New(byKeys)
	vals := btree.New(byVals)

	// Create some items.
	users := []*Item{
		&Item{Key: "user:1", Val: "Jane"},
		&Item{Key: "user:2", Val: "Andy"},
		&Item{Key: "user:3", Val: "Steve"},
		&Item{Key: "user:4", Val: "Andrea"},
		&Item{Key: "user:5", Val: "Janet"},
		&Item{Key: "user:6", Val: "Andy"},
	}

	// Insert each user into both trees
	for _, user := range users {
		keys.Set(user)
		vals.Set(user)
	}

	// Iterate over each user in the key tree
	keys.Ascend(nil, func(item interface{}) bool {
		kvi := item.(*Item)
		fmt.Printf("%s %s\n", kvi.Key, kvi.Val)
		return true
	})

	fmt.Printf("\n")
	// Iterate over each user in the val tree
	vals.Ascend(nil, func(item interface{}) bool {
		kvi := item.(*Item)
		fmt.Printf("%s %s\n", kvi.Key, kvi.Val)
		return true
	})

	// Output:
	// user:1 Jane
	// user:2 Andy
	// user:3 Steve
	// user:4 Andrea
	// user:5 Janet
	// user:6 Andy
	//
	// user:4 Andrea
	// user:2 Andy
	// user:6 Andy
	// user:1 Jane
	// user:5 Janet
	// user:3 Steve
}
```

## Performance

This implementation was designed with performance in mind. 

- `google`: The [google/btree](https://github.com/google/btree) package (without generics)
- `tidwall`: The [tidwall/btree](https://github.com/tidwall/btree) package (without generics)
- `tidwall(G)`: The [tidwall/btree](https://github.com/tidwall/btree) package (generics using the `btree.Generic` type)
- `tidwall(M)`: The [tidwall/btree](https://github.com/tidwall/btree) package (generics using the `btree.Map` type)
- `go-arr`: A simple Go array

The following benchmarks were run on my 2019 Macbook Pro (2.4 GHz 8-Core Intel Core i9) 
using Go Development version 1.18 (beta1).
The items are simple 8-byte ints. 

```
** sequential set **
google:     set-seq        1,000,000 ops in 156ms, 6,426,724/sec, 155 ns/op, 39.0 MB, 40 bytes/op
tidwall:    set-seq        1,000,000 ops in 135ms, 7,380,627/sec, 135 ns/op, 23.5 MB, 24 bytes/op
tidwall(G): set-seq        1,000,000 ops in 78ms, 12,881,995/sec, 77 ns/op, 8.2 MB, 8 bytes/op
tidwall(M): set-seq        1,000,000 ops in 46ms, 21,892,141/sec, 45 ns/op, 8.2 MB, 8 bytes/op
tidwall:    set-seq-hint   1,000,000 ops in 73ms, 13,789,017/sec, 72 ns/op, 23.5 MB, 24 bytes/op
tidwall(G): set-seq-hint   1,000,000 ops in 48ms, 20,969,431/sec, 47 ns/op, 8.2 MB, 8 bytes/op
tidwall:    load-seq       1,000,000 ops in 45ms, 22,452,523/sec, 44 ns/op, 23.5 MB, 24 bytes/op
tidwall(G): load-seq       1,000,000 ops in 22ms, 46,242,274/sec, 21 ns/op, 8.2 MB, 8 bytes/op
tidwall(M): load-seq       1,000,000 ops in 13ms, 74,371,903/sec, 13 ns/op, 8.2 MB, 8 bytes/op
go-arr:     append         1,000,000 ops in 21ms, 47,141,875/sec, 21 ns/op, 8.1 MB, 8 bytes/op

** sequential get **
google:     get-seq        1,000,000 ops in 119ms, 8,389,459/sec, 119 ns/op
tidwall:    get-seq        1,000,000 ops in 110ms, 9,068,759/sec, 110 ns/op
tidwall(G): get-seq        1,000,000 ops in 78ms, 12,813,135/sec, 78 ns/op
tidwall(M): get-seq        1,000,000 ops in 62ms, 16,053,728/sec, 62 ns/op
tidwall:    get-seq-hint   1,000,000 ops in 64ms, 15,509,696/sec, 64 ns/op
tidwall(G): get-seq-hint   1,000,000 ops in 41ms, 24,144,951/sec, 41 ns/op

** random set **
google:     set-rand       1,000,000 ops in 563ms, 1,777,592/sec, 562 ns/op, 29.7 MB, 31 bytes/op
tidwall:    set-rand       1,000,000 ops in 542ms, 1,844,397/sec, 542 ns/op, 29.6 MB, 31 bytes/op
tidwall(G): set-rand       1,000,000 ops in 234ms, 4,271,764/sec, 234 ns/op, 11.2 MB, 11 bytes/op
tidwall(M): set-rand       1,000,000 ops in 189ms, 5,292,236/sec, 188 ns/op, 11.2 MB, 11 bytes/op
tidwall:    set-rand-hint  1,000,000 ops in 602ms, 1,659,852/sec, 602 ns/op, 29.6 MB, 31 bytes/op
tidwall(G): set-rand-hint  1,000,000 ops in 278ms, 3,595,435/sec, 278 ns/op, 11.2 MB, 11 bytes/op
tidwall:    set-after-copy 1,000,000 ops in 679ms, 1,471,954/sec, 679 ns/op
tidwall(G): set-after-copy 1,000,000 ops in 238ms, 4,196,854/sec, 238 ns/op
tidwall:    load-rand      1,000,000 ops in 532ms, 1,880,877/sec, 531 ns/op, 29.6 MB, 31 bytes/op
tidwall(G): load-rand      1,000,000 ops in 232ms, 4,316,475/sec, 231 ns/op, 11.2 MB, 11 bytes/op
tidwall(M): load-rand      1,000,000 ops in 209ms, 4,790,169/sec, 208 ns/op, 11.2 MB, 11 bytes/op

** random get **
google:     get-rand       1,000,000 ops in 807ms, 1,238,703/sec, 807 ns/op
tidwall:    get-rand       1,000,000 ops in 812ms, 1,231,551/sec, 811 ns/op
tidwall(G): get-rand       1,000,000 ops in 255ms, 3,914,819/sec, 255 ns/op
tidwall(M): get-rand       1,000,000 ops in 190ms, 5,249,966/sec, 190 ns/op
tidwall:    get-rand-hint  1,000,000 ops in 876ms, 1,141,313/sec, 876 ns/op
tidwall(G): get-rand-hint  1,000,000 ops in 258ms, 3,877,775/sec, 257 ns/op

** range **
google:     ascend        1,000,000 ops in 26ms, 39,101,882/sec, 25 ns/op
tidwall:    ascend        1,000,000 ops in 20ms, 50,223,988/sec, 19 ns/op
tidwall(G): iter          1,000,000 ops in 8ms, 119,155,937/sec, 8 ns/op
tidwall(G): scan          1,000,000 ops in 6ms, 168,275,407/sec, 5 ns/op
tidwall(G): walk          1,000,000 ops in 5ms, 186,941,046/sec, 5 ns/op
go-arr:     for-loop      1,000,000 ops in 4ms, 272,234,997/sec, 3 ns/op
```

*You can find the benchmark utility at [tidwall/btree-benchmark](https://github.com/tidwall/btree-benchmark)*

## Contact

Josh Baker [@tidwall](http://twitter.com/tidwall)

## License

Source code is available under the MIT [License](/LICENSE).
