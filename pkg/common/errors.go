package common

import "errors"

func ExtendErrors(errs ...error) error {
	res := ""
	for idx, err := range errs {
		if err == nil {
			continue
		}
		res += err.Error()
		if idx != len(errs)-1 {
			res += ": "
		}
	}
	return errors.New(res)
}
