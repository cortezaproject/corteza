package rdbms

// jsonOrEmpty substitutes an empty value for a nil pointer bound to a NOT NULL
// JSON column.
//
// database/sql short-circuits a nil pointer to SQL NULL without ever calling
// Value() when the pointed-to type implements driver.Valuer, which is what the
// generated value-receiver Value() causes. The column's DEFAULT does not apply
// to an explicitly bound NULL, so the write fails the NOT NULL constraint. The
// generated queries route such columns through here instead, which is what
// their model's defaultEmptyObject already promises.
func jsonOrEmpty[T any](v *T) *T {
	if v == nil {
		return new(T)
	}

	return v
}
