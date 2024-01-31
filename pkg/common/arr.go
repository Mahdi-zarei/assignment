package common

func GetElem[T any](arr []T, block func(obj T) bool) (T, bool) {
	for _, elem := range arr {
		if block(elem) {
			return elem, true
		}
	}

	return DefaultVal[T](), false
}
