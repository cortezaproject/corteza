package ql

// MakeIdentWrapHandler makes an Ident wrapper
//
// Will return all skipped (ss) fields as-is and replace the rest with wrap, moving
// ident (Value) to args
func MakeIdentWrapHandler(wrap string, ss ...string) IdentHandler {
	return func(i Ident) (Ident, error) {
		for _, s := range ss {
			if s == i.Value {
				return i, nil
			}
		}

		i.Args = []interface{}{i.Value}
		i.Value = wrap

		return i, nil
	}
}

// MakeIdentOrderWrapHandler makes an Ident wrapper for Order columns
//
// Functionality is the same as MakeIdentWrapHandler with extra padding.
func MakeIdentOrderWrapHandler(wrap string, ss ...string) IdentHandler {
	return func(i Ident) (Ident, error) {
		for _, s := range ss {
			if s == i.Value {
				i.Value += " "
				return i, nil
			}
		}

		i.Args = []interface{}{i.Value}
		i.Value = wrap + " "

		return i, nil
	}
}
