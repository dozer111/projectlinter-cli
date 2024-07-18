package pkg

func ValueOrDefaultFromPointer[T any](val *T) T {
	var result T

	if val != nil {
		result = *val
	}

	return result
}
