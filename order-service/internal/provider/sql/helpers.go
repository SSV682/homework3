package sql

import "github.com/elgris/sqrl"

func updateIfNotNil[T any](b *sqrl.UpdateBuilder, name string, value *T) *sqrl.UpdateBuilder {
	if value != nil {
		b.Set(name, *value)
	}

	return b
}
