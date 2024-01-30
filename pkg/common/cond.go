package common

func Must1(err error) {
	if err != nil {
		panic(err)
	}
}

func Must2(_ any, err error) {
	if err != nil {
		panic(err)
	}
}

func MustGetVal[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}
