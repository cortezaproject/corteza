package id

import (
	"fmt"
	"regexp"
	"strconv"
)

type (
	ID struct {
		Num uint64
		Str [maxStrIdentifierLength]rune
	}
)

const (
	maxStrIdentifierLength = 36
)

var (
	isNumberRegex = regexp.MustCompile(`^(0|[1-9][0-9]{0,19})$`)
)

func MustNumID(n uint64) ID {
	id, err := NumID(n)
	if err != nil {
		panic(fmt.Sprintf("failed to create NumID: %v", err))
	}
	return id
}

func NumID(n uint64) (ID, error) {
	return ID{Num: n}, nil
}

func MustStrID(s string) ID {
	id, err := StrID(s)
	if err != nil {
		panic(fmt.Sprintf("failed to create StrID: %v", err))
	}
	return id
}

func StrID(s string) (ID, error) {
	if len(s) > maxStrIdentifierLength {
		return ID{}, fmt.Errorf("string too long: %s", s)
	}

	var str [maxStrIdentifierLength]rune
	for i, r := range s {
		if i < maxStrIdentifierLength {
			str[i] = r
		}
	}
	return ID{Str: str}, nil
}

func MustByteID(data []byte) (id ID) {
	id, err := ByteID(data)
	if err != nil {
		panic(fmt.Sprintf("failed to parse ID: %v", err))
	}

	return id
}

func ByteID(data []byte) (id ID, err error) {
	// numeric
	if isNumberRegex.MatchString(string(data)) {
		num, err := strconv.ParseUint(string(data), 10, 64)
		if err != nil {
			return id, err
		}
		id, err = NumID(num)
		return id, err
	}

	// check if strings fit
	if len(data) > maxStrIdentifierLength {
		return id, fmt.Errorf("string too long: %s", data)
	}

	id, err = StrID(string(data))
	return
}

func (a ID) Equal(b ID) bool {
	return a == b
}

func (id ID) IsZero() bool {
	return !id.isNum() && !id.isStr()
}

func (id ID) isNum() bool {
	return id.Num != 0
}

func (id ID) isStr() bool {
	return id.Str != [maxStrIdentifierLength]rune{}
}

func InSlice(needle ID, slice ...ID) bool {
	for _, v := range slice {
		if v.Equal(needle) {
			return true
		}
	}
	return false
}

func RemoveFromSlice(needle ID, slice ...ID) (out []ID) {
	if len(slice) == 0 {
		return
	}

	out = make([]ID, 0, len(slice)-1)

	for _, i := range slice {
		if i.Equal(needle) {
			continue
		}

		out = append(out, i)
	}

	return
}

func (id ID) Value() string {
	if id.isStr() {
		return `"` + string(id.String()) + `"`
	}

	if id.isNum() {
		return `"` + String(id.Number()) + `"`
	}

	// Backward compatibility for empty ID
	return `"0"`
}

func (id ID) String() string {
	hasValid := false
	i := 0

	for _, r := range id.Str {
		if !hasValid && r != 0 {
			hasValid = true
			continue
		}

		if hasValid && r == 0 {
			return string(id.Str[:i+1])
		}

		i++
	}

	if !hasValid {
		return ""
	}

	return string(id.Str[:])
}

func (id ID) Number() uint64 {
	return id.Num
}

// marshal/unmarshal JSON
func (id ID) MarshalJSON() ([]byte, error) {
	return []byte(id.Value()), nil
}

func (id *ID) UnmarshalJSON(data []byte) (err error) {
	// regex to check if data contains a uint64 or smaller number
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return nil
	}

	data = data[1 : len(data)-1]

	auxID, err := ByteID(data)
	if err != nil {
		return
	}

	*id = auxID
	return
}

func StringifySlice(bb ...ID) (out []string) {
	for _, b := range bb {
		out = append(out, b.Value())
	}

	return
}
