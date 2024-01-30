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
